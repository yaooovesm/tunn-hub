package authentication

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
	Tunn   *transmitter.Tunnel
}

//
// Disconnect
// @Description:
// @receiver c
//
func (c *Connection) Disconnect() {
	if c.Tunn != nil {
		_ = c.Tunn.Close()
	}
}
