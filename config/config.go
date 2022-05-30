package config

import (
	"encoding/json"
	"flag"
	log "github.com/cihub/seelog"
	"io/ioutil"
	"os"
	"tunn-hub/config/protocol"
)

var Location = ""

var Current = Config{}

//
// Config
// @Description:
//
type Config struct {
	Global      Global      `json:"global"`
	User        User        `json:"user"`
	Routes      []Route     `json:"route"`
	Device      Device      `json:"device"`
	Auth        Auth        `json:"auth"`
	DataProcess DataProcess `json:"data_process"`
	Security    Security    `json:"security"`
	Admin       Admin       `json:"admin"`
	IPPool      IPPool      `json:"ip_pool"`
	Runtime     Runtime     `json:"runtime"`
}

//
// Global
// @Description: global config
//
type Global struct {
	Tunnel
	MTU          int  `json:"mtu"`
	Pprof        int  `json:"pprof"`
	DefaultRoute bool `json:"default_route"`
	MultiConn    int  `json:"multi_connection"`
}

//
// ReadFromFile
// @Description:
// @receiver cfg
// @param path
//
func (cfg *Config) ReadFromFile(path string) {
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
	cfg.SetDefaultValue()
	_ = json.Unmarshal(bytes, cfg)
}

//
// SetDefaultValue
// @Description:
// @receiver cfg
//
func (cfg *Config) SetDefaultValue() {
	cfg.Global.Protocol = protocol.TCP
	cfg.Global.MTU = 1400
	cfg.Global.DefaultRoute = false
	//限制multi conn
	if cfg.Global.MultiConn <= 0 {
		cfg.Global.MultiConn = 1
	} else if cfg.Global.MultiConn > 32 {
		cfg.Global.MultiConn = 32
	}
}

//
// Check
// @Description:
// @receiver cfg
//
func (cfg *Config) Check() {
}

//
// MergePushed
// @Description:
// @receiver cfg
// @param push
//
func (cfg *Config) MergePushed(push PushedConfig) {
	cfg.Global.Address = push.Global.Address
	cfg.Global.Protocol = push.Global.Protocol
	cfg.Global.Port = push.Global.Port
	cfg.Global.MultiConn = push.Global.MultiConnection
	cfg.Global.MTU = push.Global.Mtu
	cfg.Routes = push.Routes
	cfg.Device = push.Device
	cfg.DataProcess.CipherType = push.DataProcess.CipherType
}

func (cfg *Config) Storage() error {
	storage := ServerConfigStorage{
		Global:      cfg.Global,
		Routes:      cfg.Routes,
		Device:      cfg.Device,
		Auth:        cfg.Auth,
		DataProcess: cfg.DataProcess,
		Security:    cfg.Security,
		Admin:       cfg.Admin,
		IPPool:      cfg.IPPool,
	}
	return storage.Dump(Location)
}

//
// Load
// @Description:
//
func Load() {
	c := flag.String("c", "", "config path")
	flag.Parse()
	Location = *c
	storage := ServerConfigStorage{}
	storage.ReadFromFile(Location)
	Current = storage.ToConfig()
	Current.Check()
	Current.Runtime.Collect()
}
