package config

import (
	"tunn-hub/config/protocol"
)

//
// Tunnel
// @Description:
//
type Tunnel struct {
	Address  string        `json:"address"`
	Port     int           `json:"port"`
	Protocol protocol.Name `json:"protocol"`
}
