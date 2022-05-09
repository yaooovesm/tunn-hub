package tcptunnel

import (
	"net"
	"tunn-hub/config"
	"tunn-hub/tunnel"
)

//
// ClientHandler
// @Description:
//
type ClientHandler struct {
}

//
// AfterInitialize
// @Description:
// @receiver h
// @param client
//
func (h ClientHandler) AfterInitialize(client *tunnel.Client) {

}

//
// CreateAndSetup
// @Description:
// @receiver h
// @param address
// @param config
// @return conn
// @return err
//
func (h *ClientHandler) CreateAndSetup(address string, config config.Config) (conn net.Conn, err error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	_ = tcpConn.SetKeepAlive(true)
	return tcpConn, nil
}
