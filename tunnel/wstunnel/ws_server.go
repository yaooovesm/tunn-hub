package wstunnel

import (
	log "github.com/cihub/seelog"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"time"
	"tunn-hub/transmitter"
	"tunn-hub/tunnel"
)

type wsonn struct {
	conn *transmitter.WSConn
	err  error
}

type ServerHandler struct {
	upgrader *websocket.Upgrader
	server   *tunnel.Server
	address  string
	ch       chan *wsonn
}

//
// AfterInitialize
// @Description:
// @receiver h
// @param server
//
func (h *ServerHandler) AfterInitialize(server *tunnel.Server) {
	h.ch = make(chan *wsonn, 1)
	h.server = server
	h.upgrader = &websocket.Upgrader{
		HandshakeTimeout: time.Second * time.Duration(45),
		CheckOrigin: func(r *http.Request) bool {
			//remoteAddr := r.RemoteAddr
			//remoteAddr = remoteAddr[0:strings.Index(remoteAddr, ":")]
			//return server.AuthServer.Check(remoteAddr)
			return true
		},
		EnableCompression: false,
	}
}

//
// CreateListener
// @Description:
// @receiver h
// @param address
// @return error
//
func (h *ServerHandler) CreateListener(address string) error {
	h.address = address
	endpoint := "/" + h.server.AuthServer.WSKey + "/access_point"
	http.HandleFunc(endpoint, func(writer http.ResponseWriter, request *http.Request) {
		wsconn, err := h.upgrader.Upgrade(writer, request, nil)
		h.ch <- &wsonn{
			conn: transmitter.WrapWSConn(wsconn),
			err:  err,
		}
	})
	go func() {
		log.Info("ws endpoint : ", endpoint)
		err := http.ListenAndServe(h.address, nil)
		if err != nil {
			_ = log.Error("ws server stopped : ", err)
		}
	}()
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
	ws := <-h.ch
	return ws.conn, ws.err
}
