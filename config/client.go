package config

//
// ClientConfig
// @Description:
//
type ClientConfig struct {
	Global      Global      `json:"global"`
	User        User        `json:"user"`
	Routes      []Route     `json:"route"`
	Device      Device      `json:"device"`
	Auth        Auth        `json:"auth"`
	DataProcess DataProcess `json:"data_process"`
	Security    Security    `json:"security"`
	Runtime     Runtime     `json:"runtime"`
	Admin       Admin       `json:"admin"`
}

//
// MergePushed
// @Description:
// @receiver cfg
// @param push
//
func (cfg *ClientConfig) MergePushed(push PushedConfig) {
	cfg.Global.Address = push.Global.Address
	cfg.Global.Protocol = push.Global.Protocol
	cfg.Global.Port = push.Global.Port
	cfg.Global.MultiConn = push.Global.MultiConnection
	cfg.Global.MTU = push.Global.Mtu
	cfg.Routes = push.Routes
	cfg.Device = push.Device
	cfg.DataProcess.CipherType = push.DataProcess.CipherType
}
