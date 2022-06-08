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
	Config config.ClientConfig
	Conn   *transmitter.Tunnel
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
