package tunnel

import (
	log "github.com/cihub/seelog"
	"net"
	"tunn-hub/administration"
	"tunn-hub/administration/model"
	"tunn-hub/authenticationv2"
	"tunn-hub/config"
	"tunn-hub/device"
	"tunn-hub/traffic"
	"tunn-hub/transmitter"
)

//
// AuthServerHandler
// @Description:
//
type AuthServerHandler struct {
	Server *Server
}

//
// GetDevice
// @Description:
// @receiver h
// @return *water.Interface
//
func (h *AuthServerHandler) GetDevice() device.Device {
	return h.Server.IFace
}

//
// AddTunnelRoute
// @Description:
// @receiver h
// @param dst
// @param tunnel
//
func (h *AuthServerHandler) AddTunnelRoute(dst string, uuid string) error {
	return h.Server.router.Add(dst, uuid)
}

//
// OnMessage
// @Description:
// @receiver h
// @param Address
// @param conn
// @param packet
//
func (h *AuthServerHandler) OnMessage(packet *authenticationv2.TransportPacket) {
	log.Info("[uuid:", packet.UUID, "] send message to server : ", string(packet.Payload))
}

//
// OnReport
// @Description:
// @receiver h
// @param packet
// @param Address
//
func (h *AuthServerHandler) OnReport(packet *authenticationv2.TransportPacket) {
	log.Info("[uuid:", packet.UUID, "] report data to server : ", len(packet.Payload), " bytes")
}

//
// AfterLogin
// @Description:
// @receiver h
// @param packet
// @param Address
//
func (h *AuthServerHandler) AfterLogin(packet *authenticationv2.TransportPacket, address string, cfg config.ClientConfig, kick func()) {
	log.Info("[Account:", cfg.User.Account, "][Address:", address, "][uuid:", packet.UUID, "] login success")
	//setup flow processor
	//tx
	txfp := traffic.NewFlowProcessor()
	//单用户TX流量统计
	txfs := &traffic.FlowStatisticsFP{Name: "tx_" + packet.UUID}
	txfp.Register(txfs, "tx_"+packet.UUID)
	//rx
	rxfp := traffic.NewFlowProcessor()
	rxfp.Name = packet.UUID
	//单用户RX流量统计
	rxfs := &traffic.FlowStatisticsFP{Name: "rx_" + packet.UUID}
	rxfp.Register(rxfs, "rx_"+packet.UUID)
	//在此处获取限速设置并注册限速器
	if cfg.Limit.Bandwidth != 0 {
		lmt := traffic.NewPPSLimiterFP(int(cfg.Limit.Bandwidth), config.Current.Global.MTU)
		txfp.Register(&lmt, "txlmt_"+packet.UUID)
		rxfp.Register(&lmt, "rxlmt_"+packet.UUID)
	}
	//在此处获取并注册流量限制器
	//单位为M
	info, _ := administration.UserServiceInstance().GetUserFullByAccount(cfg.User.Account)
	if info.Config.Limit.Flow != 0 {
		//拉取最新流量记录
		//单位转换为byte
		flowlmt := traffic.NewFlowLimitFP(txfs, info.FlowCount, uint64(info.Config.Limit.Flow)*1024*1024, kick)
		txfp.Register(&flowlmt, "flowlmt_"+packet.UUID)
		//rxfp.Register(&flowlmt, "flowlmt_"+packet.UUID)
	}
	//setup cipher
	if cfg.DataProcess.CipherType != "" {
		log.Info("[uuid:", packet.UUID, "] set cipher : ", cfg.DataProcess.CipherType)
		decryptFP := traffic.GetDecryptFP(cfg.DataProcess, cfg.DataProcess.Key)
		if decryptFP != nil {
			rxfp.Register(decryptFP, "rx_decrypt_"+packet.UUID)
		}
	}
	//setup ip table
	if cfg.Device.CIDR != "" {
		if ip, _, err := net.ParseCIDR(cfg.Device.CIDR); err == nil {
			log.Info("[", packet.UUID, "]add route record : ", ip.String()+"/32", "->", packet.UUID)
			_ = h.Server.router.Add(ip.String()+"/32", packet.UUID)
		}
	}
	//每次登录时
	h.Server.rxFlowProcessors[packet.UUID] = rxfp
	h.Server.txFlowProcessors[packet.UUID] = txfp
	//将计数器注入到multiconn中
	multiConn := transmitter.NewMultiConn(packet.UUID)
	multiConn.SetWriteFlowProcessors(txfp)
	h.Server.tunnels[packet.UUID] = multiConn
	//创建multi Conn with counter
	//注册到后台
	if administration.UserServiceInstance() != nil {
		//设置在线
		administration.UserServiceInstance().SetOnline(cfg.User.Account, address)
		//注册储存空间
		administration.UserServiceInstance().StatusService().RegisterStorage(cfg.User.Account, &model.UserStorage{
			TXFlowCounter: txfs,
			RXFlowCounter: rxfs,
			Config:        cfg,
			Address:       address,
			UUID:          packet.UUID,
		})
	}
}

//
// AfterLogout
// @Description:
// @receiver h
// @param packet
// @param Address
//
func (h *AuthServerHandler) AfterLogout(packet *authenticationv2.TransportPacket) {
	log.Info("[uuid:", packet.UUID, "] clear online")
	h.ClearOnline(packet.UUID)
}

//
// OnKick
// @Description:
// @receiver h
//
func (h *AuthServerHandler) OnKick(packet *authenticationv2.TransportPacket) {
	log.Info("[kick][uuid:", packet.UUID, "] clear online")
	h.ClearOnline(packet.UUID)
}

//
// Disconnect
// @Description:
// @receiver h
// @param Address
// @param err
//
func (h *AuthServerHandler) Disconnect(uuid string, err error) {
	log.Info("[uuid:", uuid, "] disconnected : ", err.Error())
	h.ClearOnline(uuid)
}

//
// ClearOnline
// @Description:
// @receiver h
//
func (h *AuthServerHandler) ClearOnline(uuid string) {
	if mt, ok := h.Server.tunnels[uuid]; ok && mt != nil {
		mt.Close()
		h.Server.lock.Lock()
		delete(h.Server.tunnels, uuid)
		if p, ok := h.Server.txFlowProcessors[uuid]; ok && p != nil {
			p.Close()
		}
		delete(h.Server.txFlowProcessors, uuid)
		if p, ok := h.Server.rxFlowProcessors[uuid]; ok && p != nil {
			p.Close()
		}
		delete(h.Server.rxFlowProcessors, uuid)
		h.Server.lock.Unlock()
	}
}

//
// BeforeClear
// @Description:
// @receiver h
// @param connection
//
func (h *AuthServerHandler) BeforeClear(connection *authenticationv2.Connection) {
	if connection != nil && administration.UserServiceInstance() != nil {
		administration.UserServiceInstance().SetOffline(connection.Config.User.Account)
	}
}
