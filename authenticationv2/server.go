package authenticationv2

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	log "github.com/cihub/seelog"
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
	"tunn-hub/administration"
	"tunn-hub/config"
	"tunn-hub/config/protocol"
	"tunn-hub/networking"
	"tunn-hub/transmitter"
	"tunn-hub/utils/timer"
)

//
// Server
// @Description:
//
type Server struct {
	Online        map[string]*Connection       //uuid:connection 储存在线用户
	handler       AuthServerHandler            //通过handler连接与控制服务端(hub)
	validator     IValidator                   //验证器 验证用户合法性
	IPPool        *networking.IPAddressPool    //地址池 分配IP地址
	SysRouteTable *networking.SystemRouteTable //系统路由托管
	PublicKey     []byte                       //服务器公共密钥，所有传输到服务器的数据如需要加密则用该密钥加密
	WSKey         string                       //ws&wss接入点
	version       transmitter.Version          //传输版本
	upgrader      *websocket.Upgrader          //ws&wss upgrader
}

//
// NewServer
// @Description:创建验证服务器
// @param handler
// @param validator
// @return server
// @return err
//
func NewServer(handler AuthServerHandler, validator IValidator) (server *Server, err error) {
	//检查证书文件
	if config.Current.Security.CertPem == "" || config.Current.Security.KeyPem == "" {
		return nil, ErrCertFileNotFound
	}
	//刷新密钥
	GenerateCipherKey()
	//创建实例
	s := &Server{
		Online:    make(map[string]*Connection),
		handler:   handler,
		PublicKey: config.Current.DataProcess.Key,
		version:   transmitter.V2,
	}
	//创建地址池
	poolConfig := config.Current.IPPool
	if poolConfig.Start != "" && poolConfig.End != "" && poolConfig.Network != "" {
		ipRange := networking.IPRange{}
		ipRange.Start(poolConfig.Start).End(poolConfig.End)
		iPv4AddressPool, err := networking.NewIPv4AddressPool(poolConfig.Network, ipRange)
		if err != nil {
			return nil, err
		}
		s.IPPool = iPv4AddressPool
	}
	//当通信协议为websocket时生成WSKey
	if config.Current.Global.Protocol == protocol.WS || config.Current.Global.Protocol == protocol.WSS {
		k, _ := uuid.NewV4()
		s.WSKey = hex.EncodeToString(k.Bytes())
	}
	//创建upgrader
	s.upgrader = &websocket.Upgrader{
		HandshakeTimeout: time.Second * time.Duration(20),
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: false,
	}
	//设置用户验证器
	if validator == nil {
		s.validator = DefaultValidator{}
	} else {
		s.validator = validator
	}
	//设置admin api
	administration.ServerServiceInstance().SetupAuthServer(s.KickByUUID, s.RestartByUUID, s.GetConfigByUUID, s.version, s.WSKey, s.IPPool)
	//创建完成
	return s, nil
}

//
// Start
// @Description:
// @receiver s
//
func (s *Server) Start() {
	//注册系统路由表托管
	s.SysRouteTable = networking.NewSystemRouteTable(s.handler.GetDevice().Name())
	//启动服务
	var address string
	ip := net.ParseIP(config.Current.Global.Address)
	if ip != nil {
		address = strings.Join([]string{config.Current.Auth.Address, strconv.Itoa(config.Current.Auth.Port)}, ":")
	} else {
		address = strings.Join([]string{"0.0.0.0", strconv.Itoa(config.Current.Auth.Port)}, ":")
	}
	http.HandleFunc("/authentication", func(writer http.ResponseWriter, request *http.Request) {
		wssconn, _ := s.upgrader.Upgrade(writer, request, nil)
		conn := transmitter.WrapWSConn(wssconn)
		//验证连接
		id, err := s.confirm(conn)
		if err == nil {
			go s.HandleClientConnect(conn, id)
		} else {
			log.Info("connection rejected : ", err)
		}
	})
	log.Info("authentication server listen on : ", address)
	_ = log.Error("authentication server stopped : ",
		http.ListenAndServeTLS(address, config.Current.Security.CertPem, config.Current.Security.KeyPem, nil))
}

//
// HandleClientConnect
// @Description:
// @receiver s
// @param conn
// @param uuid
//
func (s *Server) HandleClientConnect(conn net.Conn, uuid string) {
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()
	//获取UUID
	tunn := transmitter.NewTunnel(conn, s.version)
ReadLoop:
	for {
		pl, err := tunn.Read()
		if err == transmitter.ErrBadPacket {
			continue
		}
		if err != nil {
			s.DisconnectById(uuid, err)
			break ReadLoop
		}
		//组合数据包
		p := NewTransportPacket()
		err = p.Decode(pl)
		if err != nil {
			//数据包错误
			continue
		}
		if p.UUID != "" && p.Type != PacketTypeUnknown {
			switch p.Type {
			case PacketTypeLogin:
				//登录
				s.login(tunn, p, tunn.RemoteAddr().String())
			case PacketTypeLogout:
				//退出登录
				s.logout(tunn, p)
			case PacketTypeMsg:
				s.handler.OnMessage(p)
			case PacketTypeReport:
				s.handler.OnReport(p)
			}
		}
	}
}

//
// login
// @Description:
// @receiver s
// @param tunn
// @param packet
// @param address
//
func (s *Server) login(tunn *transmitter.Tunnel, packet *TransportPacket, address string) {
	cfg := config.Config{}
	err := json.Unmarshal(packet.Payload, &cfg)
	//接收用户配置
	if err != nil {
		s.reply(AuthReply{
			Ok:      false,
			Error:   "read config failed",
			Message: "读取配置失败",
		}, PacketTypeLogin, packet.UUID, tunn)
		return
	}
	//检查用户是否在线
	if s.isOnline(cfg.User.Account) {
		s.reply(AuthReply{
			Ok:      false,
			Error:   "operation failed, current user is already login",
			Message: "登录失败，当前用户已在线",
		}, PacketTypeLogin, packet.UUID, tunn)
		return
	}
	//验证用户
	clientConfig, err := s.validator.ValidateUser(cfg.User)
	if err != nil {
		s.reply(AuthReply{
			Ok:      false,
			Error:   "user authentication failed : " + err.Error(),
			Message: "用户验证失败",
		}, PacketTypeLogin, packet.UUID, tunn)
		return
	}
	//SetDeadline
	_ = tunn.SetDeadline(time.Time{})
	//set online
	pushedConfig := clientConfig.ToPushModel()
	if s.IPPool != nil {
		if pushedConfig.Device.CIDR == "" {
			ip, err := s.IPPool.DispatchCIDR(packet.UUID)
			if err == nil {
				log.Info("alloc ip address to ", packet.UUID, " : ", ip)
				pushedConfig.Device.CIDR = ip
			}
		} else {
			cidr, err := s.IPPool.StaticCIDR(packet.UUID, pushedConfig.Device.CIDR)
			if err != nil {
				log.Info("failed to alloc static ip address [", cidr, "] to ", packet.UUID, " : ", err)
				pushedConfig.Device.CIDR = ""
			} else {
				log.Info("alloc static ip address to ", packet.UUID, " : ", cidr)
				pushedConfig.Device.CIDR = cidr
			}
		}
	}
	//同步到服务端记录
	cfg.MergePushed(pushedConfig)
	//在设置路由时已经检查过冲突问题,在此处可直接应用路由
	//配置用户export路由
	if len(cfg.Routes) > 0 {
		deviceName := s.handler.GetDevice().Name()
		for i := range cfg.Routes {
			if cfg.Routes[i].Option == config.RouteOptionExport {
				log.Info("[", cfg.User.Account, "][server_dev:", deviceName, "] import route --> ", cfg.Routes[i].Network)
				//添加系统路由
				s.SysRouteTable.Merge(append([]config.Route{}, config.Route{
					Network: cfg.Routes[i].Network,
					Option:  config.RouteOptionImport,
				}))
				//networking.AddSystemRoute(cfg.Routes[i].Network, deviceName)
				//添加通道路由
				err := s.handler.AddTunnelRoute(cfg.Routes[i].Network, packet.UUID)
				if err != nil {
					_ = log.Warn("[", cfg.User.Account, "][tunnel_route:", packet.UUID, "] add route --> ",
						cfg.Routes[i].Network, " with error : ", err)
					continue
				}
				log.Info("[", cfg.User.Account, "][tunnel_route:", packet.UUID, "] add route --> ", cfg.Routes[i].Network)
			}
		}
	}
	//取消主动合入服务端暴露网络，由用户自行导入
	//合入服务端export
	//pushedConfig.Routes = append(pushedConfig.Routes, getExportRoutes()...)
	pushedConfigByte, _ := json.Marshal(pushedConfig)
	//更新系统路由
	s.SysRouteTable.DeployAll()
	//传输数据到客户端
	data := map[string]string{
		"key":     hex.EncodeToString(s.PublicKey),
		"gateway": config.Current.Device.CIDR,
		"config":  string(pushedConfigByte),
	}
	if s.WSKey != "" && len(s.WSKey) > 0 {
		data["ws_key"] = s.WSKey
	}
	b, _ := json.Marshal(data)
	s.reply(AuthReply{
		Ok:      true,
		Error:   "",
		Message: string(b),
	}, PacketTypeLogin, packet.UUID, tunn)
	s.handler.AfterLogin(packet, address, cfg)
	//设置在线状态
	s.Online[packet.UUID] = &Connection{
		UUID:   packet.UUID,
		Config: cfg,
		Conn:   tunn,
	}
}

func (s *Server) logout(tunn *transmitter.Tunnel, packet *TransportPacket) {
	//检查是否有在线
	if c, ok := s.Online[packet.UUID]; !ok || c == nil {
		s.reply(AuthReply{
			Ok:      false,
			Error:   "user not online",
			Message: "用户不在线",
		}, PacketTypeLogout, packet.UUID, tunn)
		return
	} else {
		if c.UUID != packet.UUID {
			s.reply(AuthReply{
				Ok:      false,
				Error:   "user not match",
				Message: "用户信息不匹配",
			}, PacketTypeLogout, packet.UUID, tunn)
			return
		}
	}
	cfg := &config.Config{}
	err := json.Unmarshal(packet.Payload, cfg)
	if err != nil {
		s.reply(AuthReply{
			Ok:      false,
			Error:   "read config failed",
			Message: "读取配置失败",
		}, PacketTypeLogout, packet.UUID, tunn)
		return
	}
	//验证用户
	_, err = s.validator.ValidateUser(cfg.User)
	if err != nil {
		s.reply(AuthReply{
			Ok:      false,
			Error:   "user authentication failed : " + err.Error(),
			Message: "用户验证失败",
		}, PacketTypeLogout, packet.UUID, tunn)
		return
	}
	s.reply(AuthReply{
		Ok:      true,
		Error:   "",
		Message: "登出成功",
	}, PacketTypeLogout, packet.UUID, tunn)
	//设置离线
	s.DisconnectById(packet.UUID, ErrDisconnect)
}

//
// DisconnectById
// @Description:
// @receiver s
// @param uuid
// @param err
//
func (s *Server) DisconnectById(uuid string, err error) {
	//0.当连接存在时
	if c, ok := s.Online[uuid]; ok && c != nil {
		//1.断开连接
		_ = c.Disconnect()
		//2.停止通道 (顺序：事件(上级)->Disconnect(此处)->before clear(此处)->clear(下层))
		//在进入本方法前需要处理事件
		s.handler.Disconnect(uuid, err)
		s.handler.BeforeClear(c)
		//3.返还分配的IP
		s.IPPool.ReturnBackById(uuid)
		//4.删除在线数据
		delete(s.Online, uuid)
		log.Info("[user:", c.Config.User.Account, "][uuid:", c.UUID, "] disconnected!")
	}
}

//
// confirm
// @Description:
// @receiver s
// @param conn
// @return uuid
// @return err
//
func (s *Server) confirm(conn *transmitter.WSConn) (uuid string, err error) {
	//在限时10s内，若客户端未发送已登录的UUID则连接关闭
	err = timer.TimeoutTask(func() error {
		bytes := make([]byte, 32)
		n, err := conn.Read(bytes)
		if err != nil {
			return err
		}
		uuid = string(bytes[:n])
		if uuid == "" {
			return errors.New("reject by invalid uuid")
		}
		_, err = conn.Write(bytes[:n])
		if err != nil {
			return err
		}
		return nil
	}, time.Second*10)
	if err != nil {
		return "", err
	}
	return uuid, nil
}

//
// CheckByUUID
// @Description:
// @receiver s
// @param UUID
// @return bool
//
func (s *Server) CheckByUUID(uuid string) bool {
	if c, ok := s.Online[uuid]; ok && c != nil && c.UUID == uuid {
		return true
	}
	return false
}

//
// KickByUUID
// @Description:
// @receiver s
// @param uuid
// @return error
//
func (s *Server) KickByUUID(uuid string) error {
	if c, ok := s.Online[uuid]; !ok || c == nil {
		return errors.New("user not online")
	} else {
		packet := NewTransportPacket()
		//kick
		packet.UUID = c.UUID
		s.reply(AuthReply{
			Ok:      true,
			Error:   "",
			Message: "disconnected by server",
		}, PacketTypeKick, packet.UUID, c.Conn)
		//设置离线
		s.DisconnectById(uuid, ErrKick)
		return nil
	}
}

//
// RestartByUUID
// @Description:
// @receiver s
// @param uuid
// @return error
//
func (s *Server) RestartByUUID(uuid string) error {
	if c, ok := s.Online[uuid]; !ok || c == nil {
		return errors.New("user not online")
	} else {
		packet := NewTransportPacket()
		//kick
		packet.UUID = c.UUID
		s.reply(AuthReply{
			Ok:      true,
			Error:   "",
			Message: "restart by server",
		}, PacketTypeRestart, packet.UUID, c.Conn)
		s.DisconnectById(uuid, ErrRestart)
		return nil
	}
}

//
// GetConfigByUUID
// @Description:
// @receiver s
// @param uuid
// @return cfg
// @return err
//
func (s *Server) GetConfigByUUID(uuid string) (cfg config.Config, err error) {
	if c, ok := s.Online[uuid]; ok && c != nil && c.UUID == uuid {
		return c.Config, nil
	}
	return config.Config{}, errors.New("uuid not found")
}

//
// isOnline
// @Description:
// @receiver s
// @param account
// @return bool
//
func (s *Server) isOnline(account string) bool {
	for i := range s.Online {
		if s.Online[i].Config.User.Account == account {
			return true
		}
	}
	return false
}

//
// BroadcastMsg
// @Description:
// @receiver s
// @param msg
//
func (s *Server) BroadcastMsg(msg string) {
	for id := range s.Online {
		connection := s.Online[id]
		//connection.Tunn
		s.send([]byte(msg), PacketTypeMsg, connection.UUID, connection.Conn)
	}
}

//
// SendMsgByUUID
// @Description:
// @receiver s
// @param UUID
// @param msg
//
func (s *Server) SendMsgByUUID(uuid string, msg string) {
	if connection, ok := s.Online[uuid]; ok && connection != nil {
		s.send([]byte(msg), PacketTypeMsg, connection.UUID, connection.Conn)
	}
}

//
// reply
// @Description:
// @receiver s
// @param reply
// @param t
// @param UUID
// @param conn
//
func (s *Server) reply(reply AuthReply, t PacketType, uuid string, tunn *transmitter.Tunnel) {
	defer func() {
		if err := recover(); err != nil {
			log.Info("[UUID:", uuid, "] reply error")
		}
	}()
	b, _ := json.Marshal(reply)
	p := TransportPacket{
		Type:    t,
		UUID:    uuid,
		Payload: b,
	}
	_, _ = tunn.Write(p.Encode())
}

//
// send
// @Description:
// @receiver s
// @param data
// @param t
// @param UUID
// @param conn
//
func (s *Server) send(data []byte, t PacketType, uuid string, tunn *transmitter.Tunnel) {
	defer func() {
		if err := recover(); err != nil {
			log.Info("[UUID:", uuid, "] send data failed : ", err)
		}
	}()
	p := TransportPacket{
		Type:    t,
		UUID:    uuid,
		Payload: data,
	}
	_, _ = tunn.Write(p.Encode())
}
