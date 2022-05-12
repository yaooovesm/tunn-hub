package tunnel

import (
	log "github.com/cihub/seelog"
	"net"
	"tunn-hub/administration"
	"tunn-hub/administration/model"
	"tunn-hub/authentication"
	"tunn-hub/config"
	"tunn-hub/device"
	"tunn-hub/traffic"
	"tunn-hub/transmitter"
)

//
// AuthClientHandler
// @Description:
//
type AuthClientHandler struct {
	Client *Client
}

//
// GetDevice
// @Description:
// @receiver h
// @return *water.Interface
//
func (h *AuthClientHandler) GetDevice() device.Device {
	return h.Client.IFace
}

//
// OnMessage
// @Description:
// @receiver h
// @param msg
//
func (h *AuthClientHandler) OnMessage(msg string) {
	log.Info("receive message from server : ", msg)
}

//
// OnReport
// @Description:
// @receiver h
// @param payload
//
func (h *AuthClientHandler) OnReport(payload []byte) {
	log.Info("receive report data from server : ", len(payload), " bytes")
}

//
// OnLogin
// @Description:
// @receiver h
// @param err
// @param key
//
func (h *AuthClientHandler) OnLogin(err error, key []byte) {
	if err == nil {
		//crypt
		h.Client.PK = key
		rxDecryptFP := traffic.GetDecryptFP(config.Current.DataProcess, key)
		if rxDecryptFP != nil {
			h.Client.RxFP.Register(rxDecryptFP, "rx_decrypt")
		}
		//get cipher processor
		txEncryptFP := traffic.GetEncryptFP(config.Current.DataProcess, config.Current.DataProcess.Key)
		if txEncryptFP != nil {
			h.Client.TxFP.Register(txEncryptFP, "tx_encrypt")
		}
	}
}

//
// OnLogout
// @Description:
// @receiver h
// @param err
//
func (h *AuthClientHandler) OnLogout(err error) {
	h.Client.SetErr(err)
	h.Client.Stop()
}

//
// OnDisconnect
// @Description:
// @receiver h
//
func (h *AuthClientHandler) OnDisconnect() {
	log.Info("disconnected...")
	h.Client.SetErr(ErrDisconnect)
	h.Client.Stop()
	h.Client.multiConn.Close()
}

//
// OnKick
// @Description:
// @receiver h
//
func (h *AuthClientHandler) OnKick() {
	h.Client.SetErr(ErrStoppedByServer)
	h.Client.Stop()
}

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
func (h *AuthServerHandler) OnMessage(packet *authentication.TransportPacket) {
	log.Info("[uuid:", packet.UUID, "] send message to server : ", string(packet.Payload))
}

//
// OnReport
// @Description:
// @receiver h
// @param packet
// @param Address
//
func (h *AuthServerHandler) OnReport(packet *authentication.TransportPacket) {
	log.Info("[uuid:", packet.UUID, "] report data to server : ", len(packet.Payload), " bytes")
}

//
// AfterLogin
// @Description:
// @receiver h
// @param packet
// @param Address
//
func (h *AuthServerHandler) AfterLogin(packet *authentication.TransportPacket, address string, cfg config.Config) {
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
	h.Server.txFlowCounters[packet.UUID] = txfp
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
func (h *AuthServerHandler) AfterLogout(packet *authentication.TransportPacket) {
	log.Info("[uuid:", packet.UUID, "] clear online")
	h.clearOnline(packet.UUID)
}

//
// OnKick
// @Description:
// @receiver h
//
func (h *AuthServerHandler) OnKick(packet *authentication.TransportPacket) {
	log.Info("[kick][uuid:", packet.UUID, "] clear online")
	h.clearOnline(packet.UUID)
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
	h.clearOnline(uuid)
}

//
// clearOnline
// @Description:
// @receiver h
//
func (h *AuthServerHandler) clearOnline(uuid string) {
	if mt, ok := h.Server.tunnels[uuid]; ok && mt != nil {
		mt.Close()
		h.Server.lock.Lock()
		delete(h.Server.tunnels, uuid)
		h.Server.lock.Unlock()
	}
}

//
// BeforeClear
// @Description:
// @receiver h
// @param connection
//
func (h *AuthServerHandler) BeforeClear(connection *authentication.Connection) {
	if connection != nil && administration.UserServiceInstance() != nil {
		administration.UserServiceInstance().SetOffline(connection.Config.User.Account)
	}
}
