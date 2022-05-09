package tunnel

import (
	"context"
	"errors"
	log "github.com/cihub/seelog"
	"github.com/songgao/water/waterutil"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
	"tunn-hub/administration"
	"tunn-hub/authentication"
	"tunn-hub/cache"
	"tunn-hub/common/config"
	"tunn-hub/device"
	"tunn-hub/networking"
	"tunn-hub/traffic"
	"tunn-hub/transmitter"
	"tunn-hub/utils/timer"
)

//
// Server
// @Description:
//
type Server struct {
	lock sync.Mutex
	//IFace            *water.Interface
	IFace            device.Device
	Config           config.Config
	Context          context.Context
	Cancel           context.CancelFunc
	Error            error
	AuthServer       *authentication.AuthServerV3
	tunnels          map[string]*transmitter.MultiConn
	ipTable          *cache.IpTableV2
	rxFlowProcessors map[string]*traffic.FlowProcessors
	txFlowCounters   map[string]*traffic.FlowProcessors
	TxFP             *traffic.FlowProcessors
	mtu              int
	router           *networking.RouteTable
	version          transmitter.Version
	handler          ServerConnHandler

	TXFlowCounter *traffic.FlowStatisticsFP
	RXFlowCounter *traffic.FlowStatisticsFP
}

//
// NewServer
// @Description:
// @param handler
// @return *Server
//
func NewServer(handler ServerConnHandler) *Server {
	return &Server{
		Config:  config.Current,
		mtu:     config.Current.Global.MTU,
		version: transmitter.V2,
		handler: handler,
	}
}

//
// Init
// @Description:
// @receiver s
// @return error
//
func (s *Server) Init() error {
	s.router = networking.NewRouteTable(true, 8)
	//use default
	authServerV3, err := authentication.NewServerV3(&AuthServerHandler{Server: s}, nil)
	if err != nil {
		return err
	}
	s.AuthServer = authServerV3
	ctx, cancelFunc := context.WithCancel(context.Background())
	s.Context = ctx
	s.Cancel = cancelFunc
	//iface
	if s.IFace == nil {
		dev, err := device.NewTunDevice()
		if err != nil {
			return err
		}
		err = dev.Setup()
		if err != nil {
			return err
		}
		s.IFace = dev
	}
	//cache
	s.ipTable = cache.NewIpTableV2(time.Minute*30, time.Minute*15)
	s.tunnels = make(map[string]*transmitter.MultiConn)
	s.rxFlowProcessors = make(map[string]*traffic.FlowProcessors)
	s.txFlowCounters = make(map[string]*traffic.FlowProcessors)
	s.TxFP = traffic.NewFlowProcessor()
	txEncryptFP := traffic.GetEncryptFP(config.Current.DataProcess, s.AuthServer.PublicKey)
	if txEncryptFP != nil {
		//注册tx加密
		s.TxFP.Register(txEncryptFP, "tx_encrypt")
	}
	//注册流量统计
	s.RXFlowCounter = &traffic.FlowStatisticsFP{Name: "pub_rx_flow_statistics"}
	s.TXFlowCounter = &traffic.FlowStatisticsFP{Name: "pub_tx_flow_statistics"}
	s.TxFP.Register(s.TXFlowCounter, "tx_flow_statistics")
	//注册服务器
	if administration.ServerServiceInstance() != nil {
		administration.ServerServiceInstance().SetupServer(
			s.RXFlowCounter, s.TXFlowCounter,
			s.Context, s.Cancel)
	}
	return nil
}

//
// Start
// @Description:
// @receiver s
// @return error
//
func (s *Server) Start() error {
	s.handler.AfterInitialize(s)
	go s.AuthServer.Start()
	go s.ConnectionHandler()
	go s.TXHandler()
	select {
	case <-s.Context.Done():
		err := s.Error
		s.Error = nil
		return err
	}
}

//
// Stop
// @Description:
// @receiver s
//
func (s *Server) Stop() {
	s.Cancel()
}

//
// ConnectionHandler
// @Description:
// @receiver s
//
func (s *Server) ConnectionHandler() {
	var address string
	ip := net.ParseIP(s.Config.Global.Address)
	if ip != nil {
		address = strings.Join([]string{s.Config.Global.Address, strconv.Itoa(s.Config.Global.Port)}, ":")
	} else {
		address = strings.Join([]string{"0.0.0.0", strconv.Itoa(s.Config.Global.Port)}, ":")
	}
	err := s.handler.CreateListener(address)
	if err != nil {
		s.Error = err
		s.Stop()
		return
	}
	log.Info("server listen on : ", address)
	for {
		select {
		case <-s.Context.Done():
			return
		default:
			conn, err := s.handler.AcceptConnection()
			if err != nil {
				_ = log.Warn("accept connection failed : ", err)
				_ = conn.Close()
				continue
			}
			//检查登录
			uuid, err := s.confirm(conn)
			if err != nil {
				_ = log.Warn("connection rejected : ", err)
				_ = conn.Close()
				continue
			}
			go s.RXHandler(conn, uuid)
		}
	}
}

//
// confirm
// @Description:
// @receiver s
// @param conn
// @return error
//
func (s *Server) confirm(conn net.Conn) (uuid string, err error) {
	uuid = ""
	return uuid, timer.TimeoutTask(func() error {
		bytes := make([]byte, 32)
		n, err := conn.Read(bytes)
		if err != nil {
			return err
		}
		//验证登录
		if !s.AuthServer.CheckByUUID(string(bytes[:n])) {
			return errors.New("invalid connection")
		}
		//rewrite
		_, err = conn.Write(bytes)
		if err != nil {
			return err
		}
		uuid = string(bytes)
		return nil
	}, time.Second*10)
}

//
// RXHandler
// @Description:
// @receiver s
// @param conn
//
func (s *Server) RXHandler(conn net.Conn, uuid string) {
	log.Info("rx handler start with uuid : ", uuid)
	reader := transmitter.NewTunReader(conn, s.version)
	fps, ok := s.rxFlowProcessors[uuid]
	if !ok || fps == nil {
		log.Info("[tunnel->iface] ", uuid, " 注册失败")
		return
	}
	num := -1
	s.lock.Lock()
	//注册连接
	if m, ok := s.tunnels[uuid]; ok {
		num = m.Push(conn)
	} else {
		multiConn := transmitter.NewMultiConn(uuid)
		num = multiConn.Push(conn)
		s.tunnels[uuid] = multiConn
	}
	s.lock.Unlock()
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
		log.Info("[rx][", uuid, "][#", num, "] exit")
	}()
	for {
		pl, err := reader.Read()
		if err == transmitter.ErrBadPacket {
			continue
		}
		if err != nil {
			//log.Info("[rx][", Address, "][#", num, "] exit with error ", err.Error())
			return
		}
		//从通道进入的数据需要先进行处理再分发流量
		pl = fps.Process(pl)
		//优先匹配路由
		if uuid := s.router.Route(waterutil.IPv4Destination(pl)); uuid != "" {
			if m, ok := s.tunnels[uuid]; ok {
				//重定向通道
				if txfp, ok := s.txFlowCounters[uuid]; ok {
					//计数
					_, _ = m.Get().Write(txfp.Process(s.TxFP.Process(pl)))
				} else {
					_, _ = m.Get().Write(s.TxFP.Process(pl))
				}
			}
			//路由匹配成功，无论是否发送到通道都不写入到网卡
			continue
		}
		//从数据包中生成identification
		identification := traffic.IdentificationV2(pl, traffic.In)
		if identification == "" {
			continue
		}
		//更新identification
		s.ipTable.Set(identification, uuid)
		//写入到网卡
		_, _ = s.IFace.Write(pl)
	}
}

//
// TXHandler
// @Description:
// @receiver s
//
func (s *Server) TXHandler() {
	buffer := make([]byte, s.Config.Global.MTU)
	for {
		n, err := s.IFace.Read(buffer)
		if err != nil || err == io.EOF || n == 0 {
			continue
		}
		//从数据包中生成identification
		identification := traffic.IdentificationV2(buffer, traffic.Out)
		if identification == "" {
			//未能标识流量方向做丢包处理
			continue
		}
		//匹配源方向
		if uuid, ok := s.ipTable.Get(identification); ok {
			if m, ok := s.tunnels[uuid]; ok {
				//流量TX
				//处理流量
				if txfp, ok := s.txFlowCounters[uuid]; ok {
					//计数
					_, _ = m.Get().Write(txfp.Process(s.TxFP.Process(buffer[:n])))
				} else {
					_, _ = m.Get().Write(s.TxFP.Process(buffer[:n]))
				}
			}
		} else {
			//匹配永久路由
			destination := waterutil.IPv4Destination(buffer)
			if uuid := s.router.Route(destination); uuid != "" {
				if m, ok := s.tunnels[uuid]; ok {
					//流量TX
					//处理流量
					if txfp, ok := s.txFlowCounters[uuid]; ok {
						//计数
						_, _ = m.Get().Write(txfp.Process(s.TxFP.Process(buffer[:n])))
					} else {
						_, _ = m.Get().Write(s.TxFP.Process(buffer[:n]))
					}
				}
			}
		}
	}
}
