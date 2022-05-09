package tunnel

import (
	"net"
	"tunn-hub/common/config"
)

//
// ClientConnHandler
// @Description:
//
type ClientConnHandler interface {
	//
	// AfterInitialize
	// @Description:
	// @param client
	//
	AfterInitialize(client *Client)
	//
	// CreateAndSetup
	// @Description:
	// @param Address
	// @return conn
	// @return err
	//
	CreateAndSetup(address string, config config.Config) (conn net.Conn, err error)
}

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
