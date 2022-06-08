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
	"tunn-hub/logging"
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
	flow := 2 * 1024 * 1024 * 1024
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
						break
					}
					if !bytes.Equal(pl, b) {
						fmt.Println("not equal")
						t.Fail()
						break
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
		fmt.Println("[client] transmit started...")
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

/**
100mbps -> 105%
100mbps -> 105%
50mbps  -> 105%
30mbps  -> 107%
10mbps  -> 105%
5mbps   -> 105%

*/
func TestSpeedLimiter(t *testing.T) {
	logging.Initialize()
	version := V1
	running := true
	bandwidth := 100
	random := true
	flow := 2 * bandwidth * 1024 * 1024
	size := flow / 1400
	speedExpect := float64(bandwidth) / 8 / 2
	fmt.Println("exp                = ", speedExpect, "M/s")
	fmt.Println("trans size         = ", float64(flow)/1024/1024/1024, " G")
	fmt.Println("packet total       = ", size)
	mtu := 1400
	var count int64

	var rxStart int64
	var rxEnd int64

	var txStart int64
	var txEnd int64

	// go tool pprof -http=:8000 http://127.0.0.1:9988/debug/pprof/profile
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
	//PPS
	lmt := traffic.NewPPSLimiterFP(bandwidth, mtu)
	//BPS
	//lmt := traffic.NewBPSLimiterFP(bandwidth, mtu)

	//RX
	rxfs := traffic.FlowStatisticsFP{Name: "rx"}
	rx := traffic.NewFlowProcessor()
	rx.Register(&rxfs, "rx_speed")
	rx.Register(&lmt, "tx_limiter")

	//TX
	txfs := traffic.FlowStatisticsFP{Name: "tx"}
	tx := traffic.NewFlowProcessor()
	tx.Register(&txfs, "tx_speed")
	tx.Register(&lmt, "tx_limiter")

	tcpAddr, _ := net.ResolveTCPAddr("tcp4", "0.0.0.0:8888")
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("[server] create server failed : ", err)
		return
	}
	//print
	go func() {
		for running {
			fmt.Println("[", rxfs.Name, "] packet_speed=", rxfs.PacketSpeed, "p/s flow_speed=", rxfs.FlowSpeed/1024/1024, "mb/s (", rxfs.FlowSpeed/1024, "kb/s)")
			fmt.Println("[", txfs.Name, "] packet_speed=", txfs.PacketSpeed, "p/s flow_speed=", txfs.FlowSpeed/1024/1024, "mb/s (", txfs.FlowSpeed/1024, "kb/s)")
			fmt.Println()
			fmt.Println()
			time.Sleep(time.Millisecond * 1000)
		}
	}()
	wg := sync.WaitGroup{}
	wg.Add(2)
	//RX
	maxRXLat := float64(0)
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
				fmt.Println("rx started ...")
				for {
					pl, err := reader.Read()
					if err != nil {
						fmt.Println("[server] read error : ", err)
						break
					}
					//if !bytes.Equal(pl, b) {
					//	fmt.Println("not equal")
					//	t.Fail()
					//	break
					//}
					start := time.Now().UnixNano()
					rx.Process(pl)
					end := time.Now().UnixNano()
					lat := float64(end-start) / 1000000
					if lat > maxRXLat {
						maxRXLat = lat
					}
					atomic.AddInt64(&count, 1)
					if count == int64(size) {
						break
					}
				}
				rxEnd = time.Now().UnixNano()
				wg.Done()
			}()
		}
	}()
	time.Sleep(time.Second)
	tcpAddr, _ = net.ResolveTCPAddr("tcp4", "127.0.0.1:8888")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("[client] start client failed : ", err)
		return
	}
	_ = conn.SetNoDelay(false)
	fmt.Println("[client] connected...")
	//TX
	maxTXLat := float64(0)
	go func() {
		writer := NewTunWriter(conn, version)
		txStart = time.Now().UnixNano()
		fmt.Println("tx started ...")
		for i := 0; i < size; i++ {
			if random {
				b = make([]byte, rand.Intn(mtu))
			}
			start := time.Now().UnixNano()
			data := tx.Process(b)
			end := time.Now().UnixNano()
			lat := float64(end-start) / 1000000
			if lat > maxTXLat {
				maxTXLat = lat
			}
			_, err = writer.Write(data)
			if err != nil {
				fmt.Println("[client] write error : ", err)
				return
			}
		}
		txEnd = time.Now().UnixNano()
		_ = writer.Conn().Close()
		wg.Done()
	}()
	wg.Wait()

	rxCost := rxEnd - rxStart
	rxTotal := float64(rxfs.Flow) / 1024 / 1024
	rxAvgSpeed := rxTotal / (float64(rxCost) / float64(1000000000))
	rxExpSpeed := 100 - ((rxAvgSpeed / speedExpect) * 100)
	rxAvgPacket := float64(size) / (float64(rxCost) / float64(1000000000))
	rxAvgCost := float32((float64(rxCost) / float64(size)) / 1000000)

	txCost := txEnd - txStart
	txTotal := float64(txfs.Flow) / 1024 / 1024
	txAvgSpeed := txTotal / (float64(txCost) / float64(1000000000))
	txExpSpeed := 100 - ((txAvgSpeed / speedExpect) * 100)
	txAvgPacket := float64(size) / (float64(txCost) / float64(1000000000))
	txAvgCost := float32((float64(txCost) / float64(size)) / 1000000)
	running = false
	time.Sleep(time.Millisecond * 2000)
	fmt.Println()
	fmt.Println()
	fmt.Println("传输比 ", count, "/", size)
	fmt.Println()
	fmt.Println("RX")
	fmt.Println("----------------------------------------------------------")
	fmt.Println("流量         --> ", float32(rxTotal), "M(", float32(rxTotal/1024), "G)")
	fmt.Println("数据包       --> ", rxfs.Packet)
	fmt.Println("流量平均速度  --> ", rxAvgSpeed, "M/s")
	if rxExpSpeed > 0 {
		fmt.Println("较预期速度    -->  -", rxExpSpeed, "%")
	} else {
		fmt.Println("较预期速度    -->  +", rxExpSpeed*-1, "%")
	}
	fmt.Println("包平均速度    --> ", rxAvgPacket, "pps")
	fmt.Println("平均处理时间  --> ", rxAvgCost, "ms")
	fmt.Println("最大处理时间  --> ", maxRXLat, "ms")
	fmt.Println("----------------------------------------------------------")
	fmt.Println()
	fmt.Println("TX")
	fmt.Println("----------------------------------------------------------")
	fmt.Println("流量         --> ", float32(txTotal), "M(", float32(txTotal/1024), "G)")
	fmt.Println("数据包       --> ", txfs.Packet)
	fmt.Println("流量平均速度  --> ", txAvgSpeed, "M/s")
	if txExpSpeed > 0 {
		fmt.Println("较预期速度    -->  -", txExpSpeed, "%")
	} else {
		fmt.Println("较预期速度    -->  +", txExpSpeed*-1, "%")
	}
	fmt.Println("包平均速度    --> ", txAvgPacket, "pps")
	fmt.Println("平均处理时间  --> ", txAvgCost, "ms")
	fmt.Println("最大处理时间  --> ", maxTXLat, "ms")
	fmt.Println("----------------------------------------------------------")
}
