package model

import (
	"encoding/base64"
	"encoding/json"
	"tunn-hub/config"
)

//
// ClientConfig
// @Description: 配置模型
//
type ClientConfig struct {
	Id     string         `json:"id"`
	Routes []config.Route `json:"routes"`
	Device config.Device  `json:"device"`
	Limit  config.Limit   `json:"limit"`
}

//
// Encode
// @Description:
// @receiver c
//
func (c *ClientConfig) Encode() string {
	bytes, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(bytes)
}

//
// Decode
// @Description:
// @receiver c
// @param data
// @return error
//
func (c *ClientConfig) Decode(data string) error {
	bytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}
	cfg := ClientConfig{}
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		return err
	}
	c.Id = cfg.Id
	c.Routes = cfg.Routes
	c.Device = cfg.Device
	c.Limit = cfg.Limit
	return nil
}

//
// ToStorageModel
// @Description:
// @receiver c
// @return ClientConfigStorage
//
func (c *ClientConfig) ToStorageModel() ClientConfigStorage {
	return ClientConfigStorage{
		Id:      c.Id,
		Content: c.Encode(),
	}
}

//
// ToPushModel
// @Description:
// @receiver c
// @return config.PushedConfig
//
func (c *ClientConfig) ToPushModel() config.PushedConfig {
	return config.PushedConfig{
		Global: config.PushedGlobal{
			Address:         config.Current.Global.Address,
			Port:            config.Current.Global.Port,
			Protocol:        config.Current.Global.Protocol,
			Mtu:             config.Current.Global.MTU,
			MultiConnection: config.Current.Global.MultiConn,
		},
		Routes: c.Routes,
		Device: c.Device,
		DataProcess: config.DataProcess{
			CipherType: config.Current.DataProcess.CipherType,
		},
		Limit: c.Limit,
	}
}

//
// ClientConfigStorage
// @Description: 配置储存模型
//
type ClientConfigStorage struct {
	Id      string `json:"id" gorm:"primaryKey" gorm:"column:id"`
	Content string `json:"content" gorm:"column:content"`
}

//
// TableName
// @Description:
// @receiver ClientConfigStorage
// @return string
//
func (ClientConfigStorage) TableName() string {
	return "tunn_config"
}

//
// GetConfig
// @Description:
// @receiver s
// @return cfg
// @return err
//
func (s *ClientConfigStorage) GetConfig() (cfg ClientConfig, err error) {
	cfg = ClientConfig{}
	if s.Content == "" {
		return cfg, nil
	}
	err = cfg.Decode(s.Content)
	if err != nil {
		return ClientConfig{}, err
	}
	return cfg, err
}
