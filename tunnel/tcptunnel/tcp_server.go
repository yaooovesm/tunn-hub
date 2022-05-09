package tcptunnel

import (
	"net"
	"tunn-hub/tunnel"
)

//
// ServerHandler
// @Description:
//
type ServerHandler struct {
	listener *net.TCPListener
}

//
// AfterInitialize
// @Description:
// @receiver h
// @param server
//
func (h *ServerHandler) AfterInitialize(server *tunnel.Server) {
}

//
// CreateListener
// @Description:
// @receiver h
// @param address
// @return error
//
func (h *ServerHandler) CreateListener(address string) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return err
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	h.listener = listener
	return nil
}

//
// AcceptConnection
// @Description:
// @receiver h
// @return conn
// @return err
//
func (h *ServerHandler) AcceptConnection() (conn net.Conn, err error) {
	return h.listener.AcceptTCP()
}
