package wsstunnel

import (
	"crypto/tls"
	"crypto/x509"
	log "github.com/cihub/seelog"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net"
	"net/url"
	"time"
	"tunn-hub/config"
	"tunn-hub/transmitter"
	"tunn-hub/tunnel"
)

type ClientHandler struct {
	url    url.URL
	dialer websocket.Dialer
}

//
// AfterInitialize
// @Description:
// @receiver h
// @param client
//
func (h *ClientHandler) AfterInitialize(client *tunnel.Client) {
	u := url.URL{Scheme: "wss", Host: client.Address, Path: "/" + client.AuthClient.WSKey + "/access_point"}
	h.url = u
	pool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(config.Current.Security.CertPem)
	if err != nil {
		_ = log.Error("load cert failed : ", err)
	}
	pool.AppendCertsFromPEM(ca)
	h.dialer = websocket.Dialer{
		TLSClientConfig: &tls.Config{RootCAs: pool},
		//TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		HandshakeTimeout:  time.Second * time.Duration(45),
		EnableCompression: false,
	}
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
	log.Info("connect to wss server : ", h.url.String())
	wsconn, _, err := h.dialer.Dial(h.url.String(), nil)
	return transmitter.WrapWSConn(wsconn), err
}
