package config

import "tunn-hub/config/protocol"

//
// PushedConfig
// @Description:
//
type PushedConfig struct {
	Global      PushedGlobal `json:"global"`
	Routes      []Route      `json:"route"`
	Device      Device       `json:"device"`
	DataProcess DataProcess  `json:"data_process"`
	Limit       Limit        `json:"limit"`
}

//
// PushedGlobal
// @Description:
//
type PushedGlobal struct {
	Address         string        `json:"address"`
	Port            int           `json:"port"`
	Protocol        protocol.Name `json:"protocol"`
	Mtu             int           `json:"mtu"`
	MultiConnection int           `json:"multi_connection"`
}
