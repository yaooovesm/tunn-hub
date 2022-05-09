package kcptunnel

import (
	log "github.com/cihub/seelog"
	"github.com/xtaci/kcp-go"
	"net"
	"time"
	"tunn-hub/tunnel"
)

//
// ServerHandler
// @Description:
//
type ServerHandler struct {
	listener *kcp.Listener
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
// @receiver ServerHandler
// @param address
// @return error
//
func (h *ServerHandler) CreateListener(address string) error {
	listener, err := kcp.ListenWithOptions(address, nil, 10, 3)
	h.listener = listener
	return err
}

//
// AcceptConnection
// @Description:
// @receiver ServerHandler
// @return conn
// @return err
//
func (h *ServerHandler) AcceptConnection() (conn net.Conn, err error) {
	conn, err = h.listener.AcceptKCP()
	if err != nil {
		return nil, err
	}
	err = conn.SetDeadline(time.Time{})
	if err != nil {
		log.Info("set kcp deadline failed : ", err)
	}
	return
}
