package config

import (
	"encoding/json"
	log "github.com/cihub/seelog"
	"io/ioutil"
	"os"
)

//
// ServerConfigStorage
// @Description:
//
type ServerConfigStorage struct {
	Global      Global      `json:"global"`
	Routes      []Route     `json:"route"`
	Device      Device      `json:"device"`
	Auth        Auth        `json:"auth"`
	DataProcess DataProcess `json:"data_process"`
	Security    Security    `json:"security"`
	Admin       Admin       `json:"admin"`
	IPPool      IPPool      `json:"ip_pool"`
}

//
// ReadFromFile
// @Description:
// @receiver cfg
// @param path
//
func (cfg *ServerConfigStorage) ReadFromFile(path string) {
	if path == "" {
		_ = log.Error("config not specific")
		os.Exit(-1)
		return
	}
	log.Info("load config from : ", path)
	file, err := os.OpenFile(path, os.O_RDONLY, 0600)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	if err != nil {
		_ = log.Error("failed to open config file : " + err.Error())
		os.Exit(-1)
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		_ = log.Error("failed to read config file : " + err.Error())
		os.Exit(-1)
	}
	_ = json.Unmarshal(bytes, cfg)
}

//
// ToConfig
// @Description:
// @receiver cfg
// @return Config
//
func (cfg *ServerConfigStorage) ToConfig() Config {
	return Config{
		Global:      cfg.Global,
		Routes:      cfg.Routes,
		Device:      cfg.Device,
		Auth:        cfg.Auth,
		DataProcess: cfg.DataProcess,
		Security:    cfg.Security,
		Admin:       cfg.Admin,
		IPPool:      cfg.IPPool,
	}
}

//
// Dump
// @Description:
// @receiver cfg
//
func (cfg *ServerConfigStorage) Dump(path string) error {
	bytes, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}
	log.Info("config dump to : ", path)
	return ioutil.WriteFile(path, bytes, 0600)
}
