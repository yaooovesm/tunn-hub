package authentication

import (
	"tunn-hub/config"
	"tunn-hub/device"
)

//
// AuthClientHandler
// @Description:
//
type AuthClientHandler interface {
	//
	// GetDevice
	// @Description:
	// @return *Device
	//
	GetDevice() device.Device
	//
	// OnMessage
	// @Description:
	// @param msg
	//
	OnMessage(msg string)
	//
	// OnReport
	// @Description:
	// @param payload
	//
	OnReport(payload []byte)
	//
	// OnLogin
	// @Description:
	// @param reply
	//
	OnLogin(err error, key []byte)
	//
	// OnLogout
	// @Description:
	// @param reply
	//
	OnLogout(err error)
	//
	// OnDisconnect
	// @Description:
	//
	OnDisconnect()
	//
	// OnKick
	// @Description:
	//
	OnKick()
}

//
// AuthServerHandler
// @Description:
//
type AuthServerHandler interface {
	//
	// AddTunnelRoute
	// @Description:
	// @param dst
	// @param tunnel
	// @return error
	//
	AddTunnelRoute(dst string, tunnel string) error
	//
	// GetDevice
	// @Description:
	// @return *water.Interface
	//
	GetDevice() device.Device
	//
	// OnMessage
	// @Description:
	// @param msg
	//
	OnMessage(packet *TransportPacket)
	//
	// OnReport
	// @Description:
	// @param payload
	//
	OnReport(packet *TransportPacket)
	//
	// AfterLogin
	// @Description:
	// @param reply
	//
	AfterLogin(packet *TransportPacket, address string, cfg config.Config)
	//
	// AfterLogout
	// @Description:
	// @param reply
	//
	AfterLogout(packet *TransportPacket)
	//
	// OnKick
	// @Description:
	// @param reply
	//
	OnKick(packet *TransportPacket)
	//
	// Disconnect
	// @Description:
	// @param address
	//
	Disconnect(uuid string, err error)
	//
	// BeforeClear
	// @Description:
	// @param address
	// @param connection
	//
	BeforeClear(connection *Connection)
}
