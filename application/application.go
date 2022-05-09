package application

import (
	log "github.com/cihub/seelog"
	"os"
	"tunn-hub/common/config"
	"tunn-hub/common/config/protocol"
	"tunn-hub/tunnel"
	"tunn-hub/tunnel/kcptunnel"
	"tunn-hub/tunnel/tcptunnel"
	"tunn-hub/tunnel/wsstunnel"
	"tunn-hub/tunnel/wstunnel"
	"tunn-hub/version"
)

//
// Application
// @Description:
//
type Application struct {
	Config   config.Config
	Protocol protocol.Name
}

//
// New
// @Description:
// @return *Application
//
func New() *Application {
	return &Application{
		Config:   config.Current,
		Protocol: config.Current.Global.Protocol,
	}
}

//
// Run
// @Description: run application
// @receiver app
//
func (app *Application) Run() {
	var serv = app.serverService()
	if serv == nil {
		_ = log.Warn("tunnel server type not support")
		os.Exit(-1)
		return
	}
	app.runService(serv)
}

//
// Service
// @Description:
//
type Service interface {
	//
	// Init
	// @Description: run once before start service
	// @return error
	//
	Init() error
	//
	// Start
	// @Description: start service
	// @return error
	//
	Start() error
	//
	// Stop
	// @Description: stop service
	//
	Stop()
}

//
// runService
// @Description:
// @receiver app
// @param serv
//
func (app *Application) runService(serv Service) {
	log.Info("tunnel version : ", version.Version)
	if version.Develop {
		_ = log.Warn("当前版本为测试版本！")
	}
	app.PProf()
	if serv == nil {
		_ = log.Warn("tunnel server type not support")
		os.Exit(-1)
		return
	}
	ch := make(chan error, 1)
	if err := serv.Init(); err != nil {
		_ = log.Warn("init failed : ", err)
		os.Exit(-1)
		return
	} else {
		log.Info("service init success...")
	}
	go func() {
		ch <- serv.Start()
	}()
	for {
		select {
		case err := <-ch:
			if err != nil {
				serv.Stop()
				_ = log.Warn("tunn-hub exit with error : ", err.Error())
				os.Exit(-1)
				return
			} else {
				log.Info("tunn-hub exited")
				os.Exit(0)
				return
			}
		}
	}
}

//
// serverService
// @Description:
// @receiver app
// @return Service
//
func (app *Application) serverService() Service {
	log.Info("transmit protocol : ", app.Protocol)
	switch app.Protocol {
	case protocol.TCP:
		return tunnel.NewServer(&tcptunnel.ServerHandler{})
	case protocol.KCP:
		return tunnel.NewServer(&kcptunnel.ServerHandler{})
	case protocol.WS:
		return tunnel.NewServer(&wstunnel.ServerHandler{})
	case protocol.WSS:
		return tunnel.NewServer(&wsstunnel.ServerHandler{})
	default:
		_ = log.Warn("unsupported protocol : ", app.Protocol)
	}
	return nil
}
