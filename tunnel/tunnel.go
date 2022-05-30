package tunnel

import (
	"net"
)

//
// ServerConnHandler
// @Description:
//
type ServerConnHandler interface {
	//
	// AfterInitialize
	// @Description:
	// @param server
	//
	AfterInitialize(server *Server)
	//
	// CreateListener
	// @Description:
	// @param Address
	// @return error
	//
	CreateListener(address string) error
	//
	// AcceptConnection
	// @Description:
	// @return conn
	// @return err
	//
	AcceptConnection() (conn net.Conn, err error)
}
