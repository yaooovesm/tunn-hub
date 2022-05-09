package model

import (
	"encoding/base64"
	"encoding/json"
	"tunn-hub/config"
)

//
// ClientConfig
// @Description:
//
type ClientConfig struct {
	Id     string         `json:"id"`
	Routes []config.Route `json:"routes"`
	Device config.Device  `json:"device"`
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
// ClientConfigStorage
// @Description:
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
