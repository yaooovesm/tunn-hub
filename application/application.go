package application

import (
	log "github.com/cihub/seelog"
	"os"
	"tunn-hub/config"
	"tunn-hub/config/protocol"
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
	serv     Service
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
	app.serv = app.serverService()
	StopWatcher(func() {
		app.Stop()
	})
	app.runService()
}

//
// Stop
// @Description:
// @receiver app
//
func (app *Application) Stop() {
	if app.serv == nil {
		_ = log.Error("no service running")
		return
	}
	app.serv.Stop()
	log.Info("application stopped")
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
func (app *Application) runService() {
	if app.serv == nil {
		_ = log.Error("service not support...")
		os.Exit(-1)
		return
	}
	log.Info("application version : ", version.Version)
	if version.Develop {
		_ = log.Warn("当前版本为测试版本！")
	}
	app.PProf()
	ch := make(chan error, 1)
	if err := app.serv.Init(); err != nil {
		_ = log.Warn("init failed : ", err)
		os.Exit(-1)
		return
	} else {
		log.Info("service init success...")
	}
	go func() {
		ch <- app.serv.Start()
	}()
	for {
		select {
		case err := <-ch:
			if err != nil {
				app.serv.Stop()
				_ = log.Warn("tunn-hub exit with error : ", err.Error())
				os.Exit(-1)
				return
			} else {
				log.Info("tunn-hub stopped")
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
