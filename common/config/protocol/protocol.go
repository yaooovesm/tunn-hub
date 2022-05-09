// Package protocol
// @Description: tunnel protocol
package protocol

import (
	"errors"
	"strings"
)

type Name string

var enum = map[string]Name{
	"kcp": KCP,
	"tcp": TCP,
	"udp": UDP,
	"ws":  WS,
	"wss": WSS,
}

const (
	KCP Name = "kcp"
	TCP Name = "tcp"
	UDP Name = "udp"
	WS  Name = "ws"
	WSS Name = "wss"
)

//
// FromString
// @Description:
// @param str
// @return name
// @return err
//
func FromString(str string) (name Name, err error) {
	if n, ok := enum[strings.ToLower(str)]; ok {
		return n, nil
	} else {
		return "unsupported", errors.New("unsupported protocol [" + str + "]")
	}
}
