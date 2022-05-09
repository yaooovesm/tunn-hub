package wstunnel

import (
	"bytes"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"tunn-hub/common/logging"
	"tunn-hub/traffic"
	"tunn-hub/transmitter"
)

func TestBufWSConn(t *testing.T) {
	version := transmitter.V2
	b := []byte("hello,你好！")

	upgrader := &websocket.Upgrader{
		HandshakeTimeout: time.Second * time.Duration(5),
		ReadBufferSize:   0,
		WriteBufferSize:  0,
		WriteBufferPool:  nil,
		Subprotocols:     nil,
		Error:            nil,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: true,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	s := gin.Default()
	s.GET("/5c0401661e5d46afab6348c5a2a39fa3/access_point", func(context *gin.Context) {
		wsconn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
		if err != nil {
			return
		}
		//handle conn
		go func() {
			conn := transmitter.WrapWSConn(wsconn)
			go func() {
				reader := transmitter.NewTunReader(conn, version)
				pl, err := reader.Read()
				if err != nil {
					fmt.Println("[server] read error : ", err)
					return
				}
				wg.Wait()
				fmt.Println("-----------------------------------------------------")
				fmt.Println("[server] read ", len(pl), " bytes")
				fmt.Println("string : ", string(pl))
				fmt.Println("raw    : ", pl)
				fmt.Println("raw_len: ", len(pl))
				fmt.Println("-----------------------------------------------------")
				if !bytes.Equal(pl, b) {
					fmt.Println("not equal")
					t.Fail()
				}
			}()

		}()
	})
	go func() {
		_ = s.Run("0.0.0.0:10240")
	}()

	u := url.URL{Scheme: "ws", Host: "localhost:10240", Path: "/5c0401661e5d46afab6348c5a2a39fa3/access_point"}
	fmt.Println("connecting to ", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("[client] start client failed : ", err)
		return
	}
	conn := transmitter.WrapWSConn(c)
	fmt.Println("[client] connected...")

	writer := transmitter.NewTunWriter(conn, version)
	n, err := writer.Write(b)
	if err != nil {
		fmt.Println("[client] write error : ", err)
		return
	}
	fmt.Println("-----------------------------------------------------")
	fmt.Println("[client] write ", n, " bytes")
	fmt.Println("string : ", string(b))
	fmt.Println("raw    : ", b)
	fmt.Println("raw_len: ", len(b))
	fmt.Println("-----------------------------------------------------")
	wg.Done()
	time.Sleep(time.Millisecond * 100)
}

func TestWithTransmitterSpeed(t *testing.T) {
	logging.Initialize()
	version := transmitter.V2
	flow := 10 * 1024 * 1024 * 1024
	//flow := 1024 * 10
	size := flow / 1400
	fmt.Println("trans size = ", flow/1024/1024/1024, " G")
	fmt.Println("size       = ", size)
	mtu := 1400
	var count int64

	var rxStart int64
	var rxEnd int64

	var txStart int64
	var txEnd int64

	go func() {
		log.Info("pprof serve on 0.0.0.0:9988")
		if err := http.ListenAndServe("0.0.0.0:9988", nil); err != nil {
			_ = log.Warn("pprof start failed : ", err)
		}
	}()

	var b []byte
	for i := 0; i < mtu; i++ {
		b = append(b, byte(rand.Intn(254)))
	}

	rxfs := traffic.FlowStatisticsFP{Name: "rx", Print: true}
	rx := traffic.NewFlowProcessor()
	rx.Register(&rxfs, "rx_speed")

	txfs := traffic.FlowStatisticsFP{Name: "tx", Print: true}
	tx := traffic.NewFlowProcessor()
	tx.Register(&txfs, "tx_speed")

	upgrader := &websocket.Upgrader{
		HandshakeTimeout: time.Second * time.Duration(5),
		ReadBufferSize:   0,
		WriteBufferSize:  0,
		WriteBufferPool:  nil,
		Subprotocols:     nil,
		Error:            nil,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: true,
	}

	s := gin.Default()
	s.GET("/ws", func(context *gin.Context) {
		wsconn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
		if err != nil {
			return
		}
		//handle conn
		go func() {
			conn := transmitter.WrapWSConn(wsconn)
			go func() {
				reader := transmitter.NewTunReader(conn, version)
				rxStart = time.Now().UnixNano()
				for {
					pl, err := reader.Read()
					if err != nil {
						fmt.Println("[server] read error : ", err)
						return
					}
					if !bytes.Equal(pl, b) {
						fmt.Println("not equal")
						t.Fail()
						return
					}
					rxfs.Process(pl)
					atomic.AddInt64(&count, 1)
					if count == int64(size) {
						break
					}
				}
				rxEnd = time.Now().UnixNano()
			}()

		}()
	})
	go func() {
		_ = s.Run("localhost:8090")
	}()

	u := url.URL{Scheme: "ws", Host: "localhost:8090", Path: "/ws"}
	fmt.Println("connecting to ", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("dial:", err)
	}
	conn := transmitter.WrapWSConn(c)
	fmt.Println("[client] connected...")
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		writer := transmitter.NewTunWriter(conn, version)
		txStart = time.Now().UnixNano()
		for i := 0; i < size; i++ {
			//n, err := writer.Write(txfs.Process(b))
			//fmt.Println("发送payload", n, "bytes -->", b)
			_, err := writer.Write(txfs.Process(b))
			if err != nil {
				fmt.Println("[client] write error : ", err)
				return
			}
		}
		txEnd = time.Now().UnixNano()
		wg.Done()
	}()
	wg.Wait()

	rxCost := rxEnd - rxStart
	rxTotal := float64(rxfs.Flow) / 1024 / 1024
	rxAvgSpeed := rxTotal / (float64(rxCost) / float64(1000000000))
	rxAvgPacket := float64(size) / (float64(rxCost) / float64(1000000000))
	rxAvgCost := float64(rxCost) / float64(size)

	txCost := txEnd - txStart
	txTotal := float64(txfs.Flow) / 1024 / 1024
	txAvgSpeed := txTotal / (float64(txCost) / float64(1000000000))
	txAvgPacket := float64(size) / (float64(txCost) / float64(1000000000))
	txAvgCost := float64(txCost) / float64(size)

	time.Sleep(time.Millisecond * 100)
	fmt.Println()
	fmt.Println()
	fmt.Println("----------------------------------------------------------")
	fmt.Println(count, "/", size)
	fmt.Println("rx --> packet=", rxfs.Packet, " flow =", rxTotal, "m(", rxTotal/1024, " g) avg_speed =", rxAvgSpeed, "m/s avg_cost =", rxAvgCost, "ns avg_packet_speed =", rxAvgPacket, " pps")
	fmt.Println("tx --> packet=", txfs.Packet, " flow =", txTotal, "m(", txTotal/1024, " g) avg_speed =", txAvgSpeed, "m/s avg_cost =", txAvgCost, "ns avg_packet_speed =", txAvgPacket, " pps")
	fmt.Println("----------------------------------------------------------")
}

func TestT1(t *testing.T) {
	u := url.URL{Scheme: "ws", Host: "192.168.200.101:10240", Path: "/5c0401661e5d46afab6348c5a2a39fa3/access_point"}
	fmt.Println(u.String())
	fmt.Println(u.User)
}
