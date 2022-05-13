package authentication

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	log "github.com/cihub/seelog"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"tunn-hub/config"
	"tunn-hub/config/protocol"
	"tunn-hub/transmitter"
	"tunn-hub/utils/timer"
)

//
// AuthClientV3
// @Description:
//
type AuthClientV3 struct {
	UUID      string
	sig       chan error
	handler   AuthClientHandler
	login     bool
	PublicKey []byte
	WSKey     string
	version   transmitter.Version
	tunnel    *transmitter.Tunnel
}

//
// NewClientV3
// @Description:
// @param cfg
// @return client
// @return err
//
func NewClientV3(handler AuthClientHandler) (client *AuthClientV3, err error) {
	if config.Current.Security.CertPem == "" {
		return nil, ErrCertFileNotFound
	}
	address := strings.Join([]string{config.Current.Auth.Address, strconv.Itoa(config.Current.Auth.Port)}, ":")
	log.Info("connect to authentication server : ", address)
	pool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(config.Current.Security.CertPem)
	if err != nil {
		return nil, log.Error("cannot load cert : ", err)
	}
	pool.AppendCertsFromPEM(ca)
	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{RootCAs: pool},
		//TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		HandshakeTimeout:  time.Second * time.Duration(45),
		EnableCompression: false,
	}
	u := url.URL{Scheme: "wss", Host: address, Path: "/authentication"}
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		return nil, ErrConnectFailed
	}
	wsConn := transmitter.WrapWSConn(conn)
	c := &AuthClientV3{UUID: UUID(), sig: make(chan error, 1), handler: handler, login: false, version: transmitter.V2}
	c.tunnel = transmitter.NewTunnel(wsConn, c.version)
	err = c.confirm(wsConn)
	if err != nil {
		return nil, err
	}
	go c.handle()
	return c, nil
}

//
// confirm
// @Description:
// @receiver c
// @return error
//
func (c *AuthClientV3) confirm(conn *transmitter.WSConn) error {
	return timer.TimeoutTask(func() error {
		_, err := conn.Write([]byte(c.UUID))
		if err != nil {
			return err
		}
		bytes := make([]byte, 32)
		n, err := conn.Read(bytes)
		if err != nil {
			return err
		}
		if string(bytes[:n]) != c.UUID {
			return errors.New("connection rejected")
		}
		return nil
	}, time.Second*10)
}

//
// handle
// @Description:
// @receiver c
//
func (c *AuthClientV3) handle() {
	for {
		pl, err := c.tunnel.Read()
		//pl, err := creator.Read(c.conn)
		if err == transmitter.ErrBadPacket {
			continue
		}
		if err != nil {
			c.handler.OnDisconnect()
			return
		}
		p := NewTransportPacket()
		//err = p.Decode(buffer[:n])
		err = p.Decode(pl)
		if err != nil {
			break
		}
		if p.UUID == c.UUID && p.Type != PacketTypeUnknown {
			switch p.Type {
			case PacketTypeLogin, PacketTypeLogout:
				reply := &AuthReply{}
				//解析reply
				_ = json.Unmarshal(p.Payload, reply)
				if reply.Ok {
					if p.Type == PacketTypeLogin {
						err := c.onLogin(reply)
						//login failed
						if err != nil {
							c.sig <- err
							return
						}
					}
					//logout success
					c.sig <- nil
				} else {
					c.sig <- errors.New(reply.Error)
				}
			case PacketTypeMsg:
				c.handler.OnMessage(string(p.Payload))
			case PacketTypeReport:
				c.handler.OnReport(p.Payload)
			case PacketTypeKick:
				reply := AuthReply{}
				//解析reply
				_ = json.Unmarshal(p.Payload, &reply)
				c.onKick(reply)
			}
		}
	}
}

//
// onLogin
// @Description:
// @receiver c
// @param reply
// @return error
//
func (c *AuthClientV3) onLogin(reply *AuthReply) error {
	data := make(map[string]string)
	err := json.Unmarshal([]byte(reply.Message), &data)
	if err != nil {
		return err
	}
	if config.Current.Global.DefaultRoute {
		if cidr, ok := data["gateway"]; ok && cidr != "" {
			ip, _, err := net.ParseCIDR(cidr)
			if err != nil {
				_ = log.Warn("set default_route failed : ", err)
			} else {
				cmd := exec.Command("/sbin/route", "add", "default", "gw", ip.String(), c.handler.GetDevice().Name())
				if err := cmd.Run(); err != nil {
					_ = log.Warn(" set default_route failed : ", err.Error())
				} else {
					log.Info("set default_route succeed : default gw -> ", ip.String())
				}
			}
		}
	}
	//当ws时接收ws_key
	if config.Current.Global.Protocol == protocol.WS || config.Current.Global.Protocol == protocol.WSS {
		if wskey, ok := data["ws_key"]; ok && wskey != "" {
			c.WSKey = wskey
		}
	}
	if routes, ok := data["route"]; ok && routes != "" {
		var rs []config.Route
		err := json.Unmarshal([]byte(routes), &rs)
		if err != nil {
			_ = log.Warn("failed to recv route info")
		} else {
			//merge to local
			config.Current.Routes = append(config.Current.Routes, rs...)
		}
	}
	if key, ok := data["key"]; ok && key != "" {
		keyBytes, err := hex.DecodeString(key)
		if err != nil {
			return err
		}
		c.PublicKey = keyBytes
		log.Info("receive ", len(c.PublicKey), " bytes key from server")
	}
	if cidr, ok := data["cidr"]; ok {
		log.Info("address allocated : ", cidr)
		config.Current.Device.CIDR = cidr
	} else {
		return errors.New("failed to get a address : " + data["error"])
	}
	return nil
}

//
// Logout
// @Description:
// @receiver c
// @return err
//
func (c *AuthClientV3) Logout() error {
	if !c.login {
		return nil
	}
	configBytes, err := json.Marshal(config.Current)
	if err != nil {
		return ErrAuthBadPacket
	}
	p := TransportPacket{
		Type:    PacketTypeLogout,
		UUID:    c.UUID,
		Payload: configBytes,
	}
	_, err = c.tunnel.Write(p.Encode())
	//err = p.Send(c.conn, packet.NewCreator())
	if err != nil {
		return ErrAuthBadPacket
	}
	if err := timer.TimeoutTask(func() error {
		return <-c.sig
	}, time.Second*5); err != nil {
		log.Info("failed to logout : ", err.Error())
		c.handler.OnLogout(err)
		c.handler.OnDisconnect()
		return ErrAuthFailed
	}
	log.Info("logout success")
	c.handler.OnLogout(nil)
	c.handler.OnDisconnect()
	return nil
}

//
// onKick
// @Description:
// @receiver c
// @param reply
//
func (c *AuthClientV3) onKick(reply AuthReply) {
	_ = log.Warn("server : ", reply.Message)
	c.handler.OnKick()
}

//
// Login
// @Description:
// @receiver c
// @return err
//
func (c *AuthClientV3) Login() (err error) {
	configBytes, err := json.Marshal(config.Current)
	if err != nil {
		return ErrAuthBadPacket
	}
	p := TransportPacket{
		Type:    PacketTypeLogin,
		UUID:    c.UUID,
		Payload: configBytes,
	}
	_, err = c.tunnel.Write(p.Encode())
	//err = p.Send(c.conn, packet.NewCreator())
	if err != nil {
		return err
	}
	if err := timer.TimeoutTask(func() error {
		return <-c.sig
	}, time.Second*30); err != nil {
		log.Info("failed to login : ", err.Error())
		c.handler.OnLogin(err, nil)
		if err == timer.Timeout {
			return ErrAuthTimeout
		}
		return ErrAuthFailed
	}
	log.Info("login success")
	c.handler.OnLogin(nil, c.PublicKey)
	c.login = true
	return nil
}

//
// Report
// @Description:
// @receiver c
// @param data
// @return error
//
func (c *AuthClientV3) Report(data []byte) (err error) {
	p := TransportPacket{
		Type:    PacketTypeReport,
		UUID:    c.UUID,
		Payload: data,
	}
	_, err = c.tunnel.Write(p.Encode())
	return
}

//
// Message
// @Description:
// @receiver c
// @param msg
// @return error
//
func (c *AuthClientV3) Message(msg string) (err error) {
	p := TransportPacket{
		Type:    PacketTypeMsg,
		UUID:    c.UUID,
		Payload: []byte(msg),
	}
	_, err = c.tunnel.Write(p.Encode())
	return
}
