package authenticationv2

import (
	"tunn-hub/config"
	"tunn-hub/device"
)

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
	AfterLogin(packet *TransportPacket, address string, cfg config.ClientConfig)
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
