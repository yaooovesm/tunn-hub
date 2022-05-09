package administration

import (
	"embed"
	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strconv"
	"strings"
	"tunn-hub/config"
	"tunn-hub/version"
)

/**
执行缓存优先原则，尽量减少io操作
*/

const (
	DefaultKey = "A!1b2cDe#56G*h8Z"
)

//
// ServerAdmin
// @Description:
//
type ServerAdmin struct {
	cfg           config.Admin
	db            *gorm.DB
	crypt         Crypt
	engine        *gin.Engine
	tokenServer   *TokenService
	userService   *userService
	serverService *serverService
	static        embed.FS
}

//
// NewServerAdmin
// @Description:
// @param cfg
// @return admin
// @return err
//
func NewServerAdmin(cfg config.Admin, static embed.FS) (admin *ServerAdmin, err error) {
	if cfg.DBFile == "" {
		_ = log.Error("sqlite db file not set")
		return
	}
	crypt, err := NewCrypt([]byte(DefaultKey))
	if err != nil {
		return nil, err
	}
	db, err := connectSqlite(cfg.DBFile)
	if err != nil {
		return nil, err
	}
	if !version.Develop {
		gin.SetMode(gin.ReleaseMode)
	}
	admin = &ServerAdmin{
		cfg:    cfg,
		db:     db,
		crypt:  crypt,
		engine: gin.Default(),
	}
	//token service
	admin.tokenServer = newTokenServer(crypt)
	//user service
	admin.userService = newUserService(admin)
	//server service
	admin.serverService = newServerService()
	//static
	admin.static = static
	return admin, nil
}

//
// Run
// @Description:
// @receiver ad
//
func (ad *ServerAdmin) Run() {
	//check config
	if ad.cfg.Address == "" {
		ad.cfg.Address = "0.0.0.0"
	}
	if ad.cfg.Port <= 0 || ad.cfg.Port > 65535 {
		return
	}
	var address string
	ip := net.ParseIP(ad.cfg.Address)
	if ip != nil {
		address = strings.Join([]string{ad.cfg.Address, strconv.Itoa(ad.cfg.Port)}, ":")
	} else {
		address = strings.Join([]string{"0.0.0.0", strconv.Itoa(ad.cfg.Port)}, ":")
	}
	//setup static
	ad.engine.GET("/static/*filepath", func(c *gin.Context) {
		staticServer := http.FileServer(http.FS(ad.static))
		staticServer.ServeHTTP(c.Writer, c.Request)
	})
	//redirect
	ad.engine.GET("/", func(c *gin.Context) {
		c.Request.URL.Path = "/static/"
		c.Redirect(http.StatusMovedPermanently, c.Request.URL.String())
	})
	//setup api
	api := NewServerAdminApi(ad.engine)
	api.Serv()
	//run
	go func() {
		log.Info("admin work at : ", address)
		var err error
		if ad.cfg.Https {
			log.Info("admin tls on")
			err = ad.engine.RunTLS(address, config.Current.Security.CertPem, config.Current.Security.KeyPem)
		} else {
			err = ad.engine.Run(address)
		}
		if err != nil {
			_ = log.Error("admin service stopped!")
			return
		}
	}()
}

//
// connectSqlite
// @Description:
// @param file
// @return db
// @return err
//
func connectSqlite(file string) (db *gorm.DB, err error) {
	return gorm.Open(sqlite.Open(file), &gorm.Config{})
}
