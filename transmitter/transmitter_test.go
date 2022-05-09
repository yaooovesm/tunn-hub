package transmitter

import (
	"bytes"
	"fmt"
	log "github.com/cihub/seelog"
	"math/rand"
	"net"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"tunn-hub/common/logging"
	"tunn-hub/traffic"
)

func TestReaderWriterFunction(t *testing.T) {
	version := V2
	b := []byte("hello,你好！")

	tcpAddr, _ := net.ResolveTCPAddr("tcp4", "0.0.0.0:8888")
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("[server] create server failed : ", err)
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("[server] server error :", err)
			return
		}
		reader := NewTunReader(conn, version)
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
	tcpAddr, _ = net.ResolveTCPAddr("tcp4", "127.0.0.1:8888")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("[client] start client failed : ", err)
		return
	}
	writer := NewTunWriter(conn, version)
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

func TestReaderWriterSpeed(t *testing.T) {
	logging.Initialize()
	version := V2
	flow := 10 * 1024 * 1024 * 1024
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

	tcpAddr, _ := net.ResolveTCPAddr("tcp4", "0.0.0.0:8888")
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("[server] create server failed : ", err)
		return
	}
	go func() {
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("[server] server error :", err)
				return
			}
			_ = conn.SetNoDelay(false)
			go func() {
				reader := NewTunReader(conn, version)
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
		}
	}()
	tcpAddr, _ = net.ResolveTCPAddr("tcp4", "127.0.0.1:8888")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("[client] start client failed : ", err)
		return
	}
	_ = conn.SetNoDelay(false)
	fmt.Println("[client] connected...")
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		writer := NewTunWriter(conn, version)
		txStart = time.Now().UnixNano()
		for i := 0; i < size; i++ {
			_, err = writer.Write(txfs.Process(b))
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

	/*
		v1
		----------------------------------------------------------
		1000000 / 1000000
		rx --> packet= 1000000  flow = 1335.14404296875 m avg_speed = 113.47155189837483 m/s avg_cost = 11766.3328 ns avg_packet_speed = 84988.24714527877  pps
		tx --> packet= 1000000  flow = 1335.14404296875 m avg_speed = 113.47155189837483 m/s avg_cost = 11766.3328 ns avg_packet_speed = 84988.24714527877  pps
		----------------------------------------------------------

		v2
		----------------------------------------------------------
		1000000 / 1000000
		rx --> packet= 1000000  flow = 1335.14404296875 m avg_speed = 105.45447118596059 m/s avg_cost = 12660.8576 ns avg_packet_speed = 78983.59112734986  pps
		tx --> packet= 1000000  flow = 1335.14404296875 m avg_speed = 105.45447118596059 m/s avg_cost = 12660.8576 ns avg_packet_speed = 78983.59112734986  pps
		----------------------------------------------------------
	*/
}
