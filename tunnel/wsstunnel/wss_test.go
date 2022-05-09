package wsstunnel

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"math/big"
	randm "math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"tunn-hub/logging"
	"tunn-hub/traffic"
	"tunn-hub/transmitter"
	"tunn-hub/tunnel"
)

func TestTLS(t *testing.T) {
	err := tunnel.GenerateAndSaveTlsCert("E:\\Project\\tunnel\\cmd\\x509\\")
	if err != nil {
		fmt.Println(err)
		return
	}
}

var certPath = "E:\\Project\\tunnel\\cmd\\x509\\cert.pem"
var keyPath = "E:\\Project\\tunnel\\cmd\\x509\\key.pem"

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
		EnableCompression: false,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg1 := sync.WaitGroup{}
	wg1.Add(1)
	http.HandleFunc("/5c0401661e5d46afab6348c5a2a39fa3/access_point", func(writer http.ResponseWriter, request *http.Request) {
		wsconn, err := upgrader.Upgrade(writer, request, nil)
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
		wg1.Done()
		fmt.Println(http.ListenAndServeTLS("0.0.0.0:10240", certPath, keyPath, nil))
	}()
	wg1.Wait()
	pool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(certPath)
	if err != nil {
		fmt.Println("ReadFile: ", err)
	}
	pool.AppendCertsFromPEM(ca)

	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{RootCAs: pool},
		//TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		HandshakeTimeout:  time.Second * time.Duration(45),
		EnableCompression: false,
	}

	u := url.URL{Scheme: "wss", Host: "172.18.28.101:10240", Path: "/5c0401661e5d46afab6348c5a2a39fa3/access_point"}
	fmt.Println("connecting to ", u.String())
	c, _, err := dialer.Dial(u.String(), nil)
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
	flow := 1 * 1024 * 1024 * 1024
	//flow := 10 * 1024 * 1024 * 1024
	//flow := 1400
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
		b = append(b, byte(randm.Intn(254)))
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
		EnableCompression: false,
	}

	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		wsconn, err := upgrader.Upgrade(writer, request, nil)
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
		_ = http.ListenAndServeTLS("localhost:8090", certPath, keyPath, nil)
	}()

	pool := x509.NewCertPool()

	ca, err := ioutil.ReadFile(certPath)
	if err != nil {
		fmt.Println("ReadFile: ", err)
	}
	pool.AppendCertsFromPEM(ca)

	dialer := websocket.Dialer{
		TLSClientConfig:   &tls.Config{RootCAs: pool, InsecureSkipVerify: true},
		HandshakeTimeout:  time.Second * time.Duration(45),
		EnableCompression: false,
	}

	u := url.URL{Scheme: "wss", Host: "localhost:8090", Path: "/ws"}
	fmt.Println("connecting to ", u.String())
	c, _, err := dialer.Dial(u.String(), nil)
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

func TestKeyGenerator(t *testing.T) {
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, max)
	subject := pkix.Name{ //Name代表一个X.509识别名。只包含识别名的公共属性，额外的属性被忽略。
		Organization:       []string{"go tunnel by junqirao"},
		OrganizationalUnit: []string{"junqirao"},
		CommonName:         "Go Web Programming",
	}
	template := x509.Certificate{
		SerialNumber: serialNumber, // SerialNumber 是 CA 颁布的唯一序列号，在此使用一个大随机数来代表它
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(3650 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature, //KeyUsage 与 ExtKeyUsage 用来表明该证书是用来做服务器认证的
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},               // 密钥扩展用途的序列
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}
	pk, _ := rsa.GenerateKey(rand.Reader, 2048) //生成一对具有指定字位数的RSA密钥

	//CreateCertificate基于模板创建一个新的证书
	//第二个第三个参数相同，则证书是自签名的
	//返回的切片是DER编码的证书
	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk) //DER 格式
	certOut, _ := os.Create("C:\\Users\\Administrator\\Desktop\\工作\\TLS\\cert.pem")
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICAET", Bytes: derBytes})
	certOut.Close()
	keyOut, _ := os.Create("C:\\Users\\Administrator\\Desktop\\工作\\TLS\\key.pem")
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	keyOut.Close()
}
