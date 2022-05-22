package administration

import (
	"encoding/json"
	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"tunn-hub/config"
	"tunn-hub/transmitter"
)

//
// ReportFetchRequest
// @Description:
//
type ReportFetchRequest struct {
	Token     string                       `json:"token"`
	Resources map[string]FetchResourceInfo `json:"resources"` //custom_key:resource_key
	Interval  int                          `json:"interval"`
}

//
// Decode
// @Description:
// @receiver r
// @param data
// @return error
//
func (r *ReportFetchRequest) Decode(data []byte) error {
	defer func() {
		if r.Interval <= 0 {
			r.Interval = 5000
		}
	}()
	return json.Unmarshal(data, r)
}

//
// PermissionCheck
// @Description:
// @receiver r
// @param remote
// @return error
//
func (r *ReportFetchRequest) PermissionCheck(remote net.IP) error {
	currentLevel := None
	for key := range r.Resources {
		resKey := r.Resources[key]
		if lev, ok := permissionAssignMap[resKey.Name]; ok {
			currentLevel = currentLevel.Compare(lev)
		}
	}
	return TokenServiceInstance().CheckTokenCode(r.Token, currentLevel, remote)
}

//
// Reporter
// @Description:
//
type Reporter struct {
	address  string
	engine   *gin.Engine
	upgrader *websocket.Upgrader
	clients  map[string]*websocket.Conn
	lock     sync.RWMutex
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
		clients: map[string]*websocket.Conn{},
		lock:    sync.RWMutex{},
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
		go r.HandleReportDispatch(ws)

	})
	go func() {
		log.Info("reporter work at : ", r.address)
		if config.Current.Admin.Https {
			err := http.ListenAndServeTLS(r.address, config.Current.Security.CertPem, config.Current.Security.KeyPem, nil)
			if err != nil {
				_ = log.Error("report service stopped : ", err)
			}
		} else {
			err := http.ListenAndServe(r.address, nil)
			if err != nil {
				_ = log.Error("report service stopped : ", err)
			}
		}
	}()
}

//
// HandleReportDispatch
// @Description:
// @param conn
//
func (r *Reporter) HandleReportDispatch(ws *websocket.Conn) {
	conn := transmitter.WrapWSConn(ws)
	remote := conn.RemoteAddr()
	log.Debug("[reporter] connection accepted from ", remote)
	//recv fetch request
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		_ = log.Warn("[reporter] an error occurred at read from ", remote, " : ", err)
		return
	}
	request := ReportFetchRequest{}
	err = request.Decode(buffer[:n])
	if err != nil {
		_ = log.Warn("[reporter] an error occurred at parse request from ", remote, " : ", err)
		return
	}
	remoteStr := remote.String()
	err = request.PermissionCheck(net.ParseIP(remoteStr[:strings.Index(remoteStr, ":")]))
	if err != nil {
		_ = log.Warn("[reporter] permission check failed at ", remote, " : ", err)
		return
	}
	r.lock.Lock()
	if exist, ok := r.clients[request.Token]; ok && exist != nil {
		_ = exist.Close()
	}
	r.clients[request.Token] = ws
	r.lock.Unlock()
	defer func(conn *transmitter.WSConn) {
		if err := recover(); err != nil {
			_ = log.Warn("[reporter] error from ", remote, " : ", err)
		}
		_ = ws.Close()
		delete(r.clients, request.Token)
		log.Debug("[reporter] connection from ", remote, " is closed")
	}(conn)
	stop := false
	go func() {
		n, err = conn.Read(buffer)
		if err != nil {
			return
		}
		if n > 0 {
			stop = true
		}
	}()
	for !stop {
		//处理推送
		data := FetchResources(request.Resources)
		marshal, err := json.Marshal(data)
		if err != nil {
			_, _ = conn.Write(FetchResourceResult{
				Data:  nil,
				Error: err.Error(),
			}.Byte())
			break
		}
		_, err = conn.Write(marshal)
		if err != nil {
			break
		}
		time.Sleep(time.Millisecond * time.Duration(request.Interval))
	}
}
