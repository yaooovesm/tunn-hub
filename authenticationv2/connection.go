package authenticationv2

import (
	"tunn-hub/config"
	"tunn-hub/transmitter"
)

//
// Connection
// @Description:
//
type Connection struct {
	UUID   string
	Config config.Config
	Conn   *transmitter.Tunnel
}

func (c *Connection) ReConnect() {

}

//
// Disconnect
// @Description:
// @receiver c
// @return error
//
func (c *Connection) Disconnect() error {
	return c.Conn.Close()
}
