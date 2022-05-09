package tunnel

import (
	"context"
	"errors"
	log "github.com/cihub/seelog"
	"net"
	"strconv"
	"strings"
	"time"
	"tunn-hub/authentication"
	"tunn-hub/config"
	"tunn-hub/device"
	"tunn-hub/traffic"
	"tunn-hub/transmitter"
	"tunn-hub/utils/timer"
)

//
// Client
// @Description:
//
type Client struct {
	//IFace            *water.Interface
	IFace            device.Device
	Config           config.Config
	Context          context.Context
	Cancel           context.CancelFunc
	Error            error
	AuthClient       *authentication.AuthClientV3
	Address          string
	tunnelIndex      int
	maxIndex         int
	size             int
	multiConn        *transmitter.MultiConn
	TxFP             *traffic.FlowProcessors
	RxFP             *traffic.FlowProcessors
	mtu              int
	PK               []byte
	version          transmitter.Version
	handler          ClientConnHandler
	running          bool
	txHandlerRunning bool
}

//
// NewClient
// @Description: 持久化多连接接模型
// @return *TCPClientV3
//
func NewClient(handler ClientConnHandler) *Client {
	size := config.Current.Global.MultiConn
	if size <= 0 {
		size = 1
	}
	if size > 32 {
		size = 32
	}
	log.Info("multi_connection set to ", size)
	return &Client{
		IFace:            nil,
		Config:           config.Current,
		Error:            nil,
		size:             size,
		tunnelIndex:      0,
		maxIndex:         size - 1, // size -1
		multiConn:        nil,
		mtu:              config.Current.Global.MTU,
		version:          transmitter.V2,
		handler:          handler,
		running:          false,
		txHandlerRunning: false,
	}
}

//
// Init
// @Description:
// @receiver C
// @return error
//
func (c *Client) Init() error {
	//preprocess Address
	c.Address = strings.Join([]string{c.Config.Global.Address, strconv.Itoa(c.Config.Global.Port)}, ":")
	//rx flow processor
	c.RxFP = traffic.NewFlowProcessor()
	c.RxFP.Name = "client_rx"
	//"RX : rx_packet_speed=", TXFs.PacketSpeed, "p/s rx_flow_speed=", TXFs.FlowSpeed/1024/1024, "mb/s"
	RXFs := &traffic.FlowStatisticsFP{Name: "rx"}
	c.RxFP.Register(RXFs, "rx_fs")
	//tx flow processor
	c.TxFP = traffic.NewFlowProcessor()
	c.TxFP.Name = "client_tx"
	//"TX : tx_packet_speed=", TXFs.PacketSpeed, "p/s tx_flow_speed=", TXFs.FlowSpeed/1024/1024, "mb/s"
	TXFs := &traffic.FlowStatisticsFP{Name: "tx"}
	c.TxFP.Register(TXFs, "tx_fs")
	return nil
}

//
// Start
// @Description:
// @receiver C
// @return error
//
func (c *Client) Start() error {
	//update key
	config.GenerateCipherKey()
	//multi conn
	c.multiConn = transmitter.NewMultiConn("client")
	//context
	ctx, cancelFunc := context.WithCancel(context.Background())
	c.Context = ctx
	c.Cancel = cancelFunc
	//auth
	//client, err := authentication.NewClientV2(c.Config.Auth, &AuthClientHandler{Client: c})
	clientV3, err := authentication.NewClientV3(&AuthClientHandler{Client: c})
	if err != nil {
		return err
	}
	c.AuthClient = clientV3
	//login
	err = c.AuthClient.Login()
	if err != nil {
		return err
	}
	//get interface cidr address after login
	//iface
	if c.IFace == nil {
		dev, err := device.NewTunDevice()
		if err != nil {
			return err
		}
		err = dev.Setup()
		if err != nil {
			return err
		}
		c.IFace = dev
	} else {
		//客户端可能重置网卡IP
		err := c.IFace.OverwriteCIDR(config.Current.Device.CIDR)
		if err != nil {
			return err
		}
	}
	c.handler.AfterInitialize(c)
	for i := 0; i < c.size; i++ {
		conn, err := c.handler.CreateAndSetup(c.Address, c.Config)
		if err != nil {
			return err
		}
		//confirm
		err = c.confirm(conn)
		if err != nil {
			return err
		}
		num := c.multiConn.Push(conn)
		go c.RXHandler(conn, num)
	}
	c.running = true
	if !c.txHandlerRunning {
		go c.TXHandler()
	}
	select {
	case <-c.Context.Done():
		c.running = false
		err := c.Error
		c.Error = nil
		return err
	}
}

//
// confirm
// @Description: 验证通道UUID
// @receiver c
// @param conn
//
func (c *Client) confirm(conn net.Conn) error {
	return timer.TimeoutTask(func() error {
		uuid := c.AuthClient.UUID
		_, err := conn.Write([]byte(uuid))
		if err != nil {
			return err
		}
		bytes := make([]byte, 32)
		n, err := conn.Read(bytes)
		if err != nil {
			return err
		}
		if uuid != string(bytes[:n]) {
			return errors.New("connection rejected")
		}
		return nil
	}, time.Second*10)
}

//
// Stop
// @Description:
// @receiver C
//
func (c *Client) Stop() {
	c.Cancel()
}

//
// SetErr
// @Description:
// @receiver c
// @param err
//
func (c *Client) SetErr(err error) {
	c.Error = err
}

//
// Logout
// @Description:
// @receiver c
//
func (c *Client) Logout() {
	defer func() {
		if c.multiConn != nil {
			c.multiConn.Close()
		}
	}()
	err := c.AuthClient.Logout()
	if err != nil {
		_ = log.Warn("logout failed : ", err.Error())
		return
	}
	log.Info("client application logout")
}

//
// TXHandler
// @Description:
// @receiver c
//
func (c *Client) TXHandler() {
	defer func() {
		log.Info("[tx_handler] exit")
		c.txHandlerRunning = false
	}()
	c.txHandlerRunning = true
	packet := make([]byte, c.Config.Global.MTU)
	for {
		//interface -> tunnel
		n, err := c.IFace.Read(packet)
		if err != nil || n == 0 {
			continue
		}
		//流量处理
		if c.running {
			if conn := c.multiConn.Get(); conn != nil {
				//先处理流量再封包
				_, _ = conn.Write(c.TxFP.Process(packet[:n]))
			}
		}
	}
}

//
// RXHandler
// @Description:
// @receiver C
//
func (c *Client) RXHandler(conn net.Conn, num int) {
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()
	//封包器
	reader := transmitter.NewTunReader(conn, c.version)
	for {
		pl, err := reader.Read()
		if err != nil && err != transmitter.ErrBadPacket {
			log.Info("[rx][#", num, "] exit with error ", err.Error())
			c.SetErr(ErrDisconnectAccidentally)
			c.Stop()
			return
		}
		//流量处理
		_, _ = c.IFace.Write(c.RxFP.Process(pl))
	}
}
