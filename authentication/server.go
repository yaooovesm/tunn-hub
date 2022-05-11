package authentication

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
	"sync"
	"time"
	"tunn-hub/administration"
	"tunn-hub/config"
	"tunn-hub/config/protocol"
	"tunn-hub/networking"
	"tunn-hub/transmitter"
	"tunn-hub/utils/timer"
)

//
// AuthServerV3
// @Description:
//
type AuthServerV3 struct {
	online        map[string]*Connection
	handler       AuthServerHandler
	lock          sync.Mutex
	PublicKey     []byte
	WSKey         string
	version       transmitter.Version
	upgrader      *websocket.Upgrader
	validator     IValidator
	IPPool        *networking.IPAddressPool
	SysRouteTable *networking.SystemRouteTable
}

//
// NewServerV3
// @Description:
// @return *AuthServer
//
func NewServerV3(handler AuthServerHandler, validator IValidator) (server *AuthServerV3, err error) {
	//生成32位公共秘钥
	GenerateCipherKey()
	//在客户端登录成功时通过report包发送到客户端
	s := &AuthServerV3{
		online:    make(map[string]*Connection),
		handler:   handler,
		PublicKey: config.Current.DataProcess.Key,
		version:   transmitter.V2,
	}
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
	if config.Current.Global.Protocol == protocol.WS || config.Current.Global.Protocol == protocol.WSS {
		//当通信协议为websocket时生成WSKey
		k, _ := uuid.NewV4()
		s.WSKey = hex.EncodeToString(k.Bytes())
	}
	if config.Current.Security.CertPem == "" || config.Current.Security.KeyPem == "" {
		return nil, ErrCertFileNotFound
	}
	s.upgrader = &websocket.Upgrader{
		HandshakeTimeout: time.Second * time.Duration(20),
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: false,
	}
	if validator == nil {
		s.validator = DefaultValidator{}
	} else {
		s.validator = validator
	}
	administration.ServerServiceInstance().SetupAuthServer(s.KickByUUID, s.GetConfigByUUID, s.version, s.WSKey, s.IPPool)
	return s, nil
}

//
// Start
// @Description:
// @receiver s
//
func (s *AuthServerV3) Start() {
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
// confirm
// @Description:
// @receiver s
// @param conn
// @return uuid
// @return err
//
func (s *AuthServerV3) confirm(conn *transmitter.WSConn) (uuid string, err error) {
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
func (s *AuthServerV3) CheckByUUID(uuid string) bool {
	for i := range s.online {
		if s.online[i].UUID == uuid {
			return true
		}
	}
	return false
}

//
// GetConfigByUUID
// @Description:
// @receiver s
// @param uuid
// @return config.Config
//
func (s *AuthServerV3) GetConfigByUUID(uuid string) (cfg config.Config, err error) {
	for i := range s.online {
		if c, ok := s.online[i]; ok && c.UUID == uuid {
			return c.Config, nil
		}
	}
	return config.Config{}, errors.New("uuid not found")
}

//
// HandleClientConnect
// @Description:
// @receiver s
// @param conn
//
func (s *AuthServerV3) HandleClientConnect(conn net.Conn, uuid string) {
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()
	//获取UUID
	tunn := transmitter.NewTunnel(conn, s.version)
	for {
		pl, err := tunn.Read()
		//pl, err := creator.Read(conn)
		if err == transmitter.ErrBadPacket {
			continue
		}
		if err != nil {
			s.handler.Disconnect(uuid, err)
			s.clearByUUID(uuid)
			return
		}
		//组合数据包
		p := NewTransportPacket()
		//err = p.Decode(buffer[:n])
		err = p.Decode(pl)
		if err != nil {
			continue
		}
		if p.UUID != "" && p.Type != PacketTypeUnknown {
			switch p.Type {
			case PacketTypeLogin:
				s.login(tunn, p, GetRemoteAddr(conn))
			case PacketTypeLogout:
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
// isOnline
// @Description:
// @receiver s
// @param account
// @return bool
//
func (s *AuthServerV3) isOnline(account string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	for i := range s.online {
		if s.online[i].Config.User.Account == account {
			return true
		}
	}
	return false
}

//
// login
// @Description:
// @param conn
//
func (s *AuthServerV3) login(tunn *transmitter.Tunnel, packet *TransportPacket, address string) {
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
	s.lock.Lock()
	pushedConfig := clientConfig.ToPushModel()
	if s.IPPool != nil {
		if pushedConfig.Device.CIDR == "" {
			ip, err := s.IPPool.DispatchCIDR(packet.UUID)
			if err == nil {
				log.Info("alloc ip address to ", packet.UUID, " : ", ip)
				pushedConfig.Device.CIDR = ip
			}
		} else {
			cidr, _ := s.IPPool.PickCIDR(pushedConfig.Device.CIDR, packet.UUID)
			pushedConfig.Device.CIDR = cidr
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
	//合入服务端export
	pushedConfig.Routes = append(pushedConfig.Routes, getExportRoutes()...)
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
	s.online[packet.UUID] = &Connection{
		UUID:   packet.UUID,
		Config: cfg,
		Tunn:   tunn,
	}
	s.lock.Unlock()
}

//
// login
// @Description:
// @param conn
//
func (s *AuthServerV3) logout(tunn *transmitter.Tunnel, packet *TransportPacket) {
	//检查是否有在线
	s.lock.Lock()
	if c, ok := s.online[packet.UUID]; !ok || c == nil {
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
	s.lock.Unlock()
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
		}, PacketTypeLogin, packet.UUID, tunn)
		return
	}
	s.clearByUUID(packet.UUID)
	log.Info("[authentication][user:", cfg.User.Account, "] logout success")
	s.reply(AuthReply{
		Ok:      true,
		Error:   "",
		Message: "登出成功",
	}, PacketTypeLogout, packet.UUID, tunn)
	s.handler.AfterLogout(packet)
}

//
// KickByUUID
// @Description:
// @receiver s
// @param uuid
// @return error
//
func (s *AuthServerV3) KickByUUID(uuid string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if c, ok := s.online[uuid]; !ok || c == nil {
		return errors.New("user not online")
	} else {
		packet := NewTransportPacket()
		//kick
		packet.UUID = c.UUID
		s.reply(AuthReply{
			Ok:      true,
			Error:   "",
			Message: "disconnected by server",
		}, PacketTypeKick, packet.UUID, c.Tunn)
		go func() {
			log.Info("[uuid:", uuid, "] connection will be clean in 10s")
			time.Sleep(time.Second * 10)
			//clear
			s.handler.OnKick(packet)
			s.handler.BeforeClear(s.online[uuid])
			delete(s.online, uuid)
		}()
		return nil
	}
}

//
// clearByUUID
// @Description:
// @receiver s
// @param uuid
//
func (s *AuthServerV3) clearByUUID(uuid string) {
	s.handler.BeforeClear(s.online[uuid])
	s.lock.Lock()
	if c, ok := s.online[uuid]; ok && c != nil {
		cidr := c.Config.Device.CIDR
		if cidr != "" {
			ip, _, err := net.ParseCIDR(cidr)
			if err == nil {
				s.IPPool.ReturnBack(ip.To4().String())
			}
		}
	}
	delete(s.online, uuid)
	s.lock.Unlock()
}

//
// BroadcastMsg
// @Description:
// @receiver s
// @param msg
//
func (s *AuthServerV3) BroadcastMsg(msg string) {
	for id := range s.online {
		connection := s.online[id]
		//connection.Tunn
		s.send([]byte(msg), PacketTypeMsg, connection.UUID, connection.Tunn)
	}
}

//
// SendMsgByUUID
// @Description:
// @receiver s
// @param UUID
// @param msg
//
func (s *AuthServerV3) SendMsgByUUID(uuid string, msg string) {
	if connection, ok := s.online[uuid]; ok && connection != nil {
		s.send([]byte(msg), PacketTypeMsg, connection.UUID, connection.Tunn)
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
func (s *AuthServerV3) reply(reply AuthReply, t PacketType, uuid string, tunn *transmitter.Tunnel) {
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
func (s *AuthServerV3) send(data []byte, t PacketType, uuid string, tunn *transmitter.Tunnel) {
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
	//_ = p.Send(conn, packet.NewCreator())
}
