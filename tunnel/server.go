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
	"tunn-hub/authenticationv2"
	"tunn-hub/cache"
	"tunn-hub/config"
	"tunn-hub/device"
	"tunn-hub/networking"
	"tunn-hub/traffic"
	"tunn-hub/transmitter"
	"tunn-hub/utils/timer"
)

const (
	RXCounterDumpFile = ".rx_counter"
	TXCounterDumpFile = ".tx_counter"
)

//
// Server
// @Description:
//
type Server struct {
	lock             sync.Mutex
	IFace            device.Device
	Config           config.Config
	Context          context.Context
	Cancel           context.CancelFunc
	Error            error
	AuthServer       *authenticationv2.Server
	tunnels          map[string]*transmitter.MultiConn
	ipTable          *cache.IpTableV2
	rxFlowProcessors map[string]*traffic.FlowProcessors
	txFlowProcessors map[string]*traffic.FlowProcessors
	TxFP             *traffic.FlowProcessors
	RxFP             *traffic.FlowProcessors
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
	//authServerV3, err := authentication.NewServerV3(&AuthServerHandler{Server: s}, nil)
	//if err != nil {
	//	return err
	//}
	//s.AuthServer = authServerV3
	authServer, err := authenticationv2.NewServer(&AuthServerHandler{Server: s}, nil)
	if err != nil {
		return err
	}
	s.AuthServer = authServer
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
	s.txFlowProcessors = make(map[string]*traffic.FlowProcessors)
	//??????RX????????????
	s.RxFP = traffic.NewFlowProcessor()
	s.RxFP.Name = "global_rx"
	//??????????????????
	rxFlowCounter := &traffic.FlowStatisticsFP{Name: "pub_rx_flow_statistics"}
	s.RxFP.Register(rxFlowCounter, "rx_flow_statistics")
	//???????????????????????????,????????????????????????
	err = rxFlowCounter.LoadFromDump(RXCounterDumpFile)
	if err != nil {
		_ = log.Warn("failed to load rx count record : ", err)
	}
	//??????TX????????????
	s.TxFP = traffic.NewFlowProcessor()
	s.TxFP.Name = "global_tx"
	txEncryptFP := traffic.GetEncryptFP(config.Current.DataProcess, s.AuthServer.PublicKey)
	if txEncryptFP != nil {
		//??????tx??????
		s.TxFP.Register(txEncryptFP, "tx_encrypt")
	}
	txFlowCounter := &traffic.FlowStatisticsFP{Name: "pub_tx_flow_statistics"}
	s.TxFP.Register(txFlowCounter, "tx_flow_statistics")
	//???????????????????????????,????????????????????????
	err = txFlowCounter.LoadFromDump(TXCounterDumpFile)
	if err != nil {
		_ = log.Warn("failed to load tx count record : ", err)
	}
	//???????????????
	if administration.ServerServiceInstance() != nil {
		administration.ServerServiceInstance().SetupServer(
			rxFlowCounter, txFlowCounter,
			s.Context, s.Cancel)
	}
	s.RXFlowCounter = rxFlowCounter
	s.TXFlowCounter = txFlowCounter
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
	//????????????????????????
	err := s.RXFlowCounter.Dump(RXCounterDumpFile)
	if err != nil {
		_ = log.Warn("failed to dump rx counter to file : ", RXCounterDumpFile)
	} else {
		log.Info("rx counter dumped")
	}
	err = s.TXFlowCounter.Dump(TXCounterDumpFile)
	if err != nil {
		_ = log.Warn("failed to dump tx counter to file : ", TXCounterDumpFile)
	} else {
		log.Info("tx counter dumped")
	}
	//??????????????????
	log.Info("commit all users flow counter")
	administration.UserServiceInstance().CommitAllFlowCount()
	log.Info("commit done")
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
			//????????????
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
		//????????????
		recv := string(bytes[:n])
		if !s.AuthServer.CheckByUUID(recv) {
			log.Info("[uuid:", recv, "] rejected")
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
		log.Info("[tunnel->iface] ", uuid, " ????????????")
		return
	}
	num := -1
	s.lock.Lock()
	//????????????
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
		//????????????
		s.RxFP.Process(pl)
		//????????????????????????????????????????????????????????????
		pl = fps.Process(pl)
		//??????????????????
		if uuid := s.router.Route(waterutil.IPv4Destination(pl)); uuid != "" {
			if m, ok := s.tunnels[uuid]; ok {
				//???????????????
				_, _ = m.Write(s.TxFP.Process(pl))
			}
			//?????????????????????????????????????????????????????????????????????
			continue
		}
		//?????????????????????identification
		identification := traffic.IdentificationV2(pl, traffic.In)
		if identification == "" {
			continue
		}
		//??????identification
		s.ipTable.Set(identification, uuid)
		//???????????????
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
		pl := buffer[:n]
		//?????????????????????identification
		identification := traffic.IdentificationV2(pl, traffic.Out)
		if identification == "" {
			//???????????????????????????????????????
			continue
		}

		//???????????????
		if uuid, ok := s.ipTable.Get(identification); ok {
			if m, ok := s.tunnels[uuid]; ok {
				//??????TX
				//????????????
				_, _ = m.Write(s.TxFP.Process(pl))
				continue
			}
		}
		//????????????
		destination := waterutil.IPv4Destination(pl)
		if uuid := s.router.Route(destination); uuid != "" {
			if m, ok := s.tunnels[uuid]; ok {
				//??????TX
				//????????????
				_, _ = m.Write(s.TxFP.Process(pl))
			}
		}
	}
}
