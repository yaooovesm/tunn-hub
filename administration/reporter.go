package administration

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
	"tunn-hub/config"
	"tunn-hub/transmitter"
)

type Reporter struct {
	address  string
	engine   *gin.Engine
	upgrader *websocket.Upgrader
}

//
// NewReporter
// @Description:
// @param cfg
// @return *Reporter
//
func NewReporter(cfg config.Admin) *Reporter {
	var address string
	ip := net.ParseIP(cfg.Address)
	if ip != nil {
		address = strings.Join([]string{cfg.Address, strconv.Itoa(cfg.ReporterPort)}, ":")
	} else {
		address = strings.Join([]string{"0.0.0.0", strconv.Itoa(cfg.ReporterPort)}, ":")
	}
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	return &Reporter{
		address: address,
		engine:  engine,
		upgrader: &websocket.Upgrader{
			HandshakeTimeout: time.Second * time.Duration(45),
			CheckOrigin: func(r *http.Request) bool {
				//remoteAddr := r.RemoteAddr
				//remoteAddr = remoteAddr[0:strings.Index(remoteAddr, ":")]
				//return server.AuthServer.Check(remoteAddr)
				return true
			},
			EnableCompression: false,
		},
	}
}

//
// Serve
// @Description:
// @receiver r
//
func (r *Reporter) Serve() {
	http.HandleFunc("/reporter", func(writer http.ResponseWriter, request *http.Request) {
		ws, err := r.upgrader.Upgrade(writer, request, nil)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		//err = TokenServiceInstance().CheckToken(ctx, User)
		//if err != nil {
		//	_, _ = writer.Write([]byte(err.Error()))
		//	return
		//}
		go HandleWSApiDispatch(ws)

	})
	go func() {
		log.Info("reporter work at : ", r.address)
		err := http.ListenAndServe(r.address, nil)
		if err != nil {
			_ = log.Error("report service stopped : ", err)
		}
	}()
}

//
// HandleWSApiDispatch
// @Description:
// @param conn
//
func HandleWSApiDispatch(ws *websocket.Conn) {
	conn := transmitter.WrapWSConn(ws)
	defer func(conn *transmitter.WSConn) {
		_ = conn.Close()
	}(conn)
	log.Info("[reporter] connection accepted from ", conn.RemoteAddr())
	buffer := make([]byte, 4096)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			_ = log.Warn("[reporter] an error occurred at ", conn.RemoteAddr(), " : ", err)
			break
		}
		fmt.Println("ws read --> ", string(buffer[:n]))
	}
}
