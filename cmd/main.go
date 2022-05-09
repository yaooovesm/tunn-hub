package main

import (
	"embed"
	log "github.com/cihub/seelog"
	"tunn-hub/administration"
	"tunn-hub/application"
	"tunn-hub/config"
	"tunn-hub/logging"
)

//go:embed static
var static embed.FS

//
// main
// @Description: entrance
//
func main() {
	//initialize log
	logging.Initialize()
	//load config
	config.Load()
	//create application
	app := application.New()
	//create admin
	admin, err := administration.NewServerAdmin(config.Current.Admin, static)
	if err != nil {
		log.Info("create admin failed : ", err)
		return
	}
	admin.Run()
	//run application
	app.Run()
}
