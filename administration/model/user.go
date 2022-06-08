package model

import (
	"encoding/json"
	"tunn-hub/config"
	"tunn-hub/traffic"
)

//
// User
// @Description:
//
type User struct {
	UserInfo   `json:"info"`
	UserStatus `json:"status"`
}

//
// UserInfo
// @Description: from database
//
type UserInfo struct {
	Id         string `json:"id" gorm:"primaryKey" gorm:"column:id"`
	Account    string `json:"account" gorm:"column:account"`
	Password   string `json:"password" gorm:"column:password"`
	Email      string `json:"email" gorm:"column:email"`
	Auth       string `json:"auth" gorm:"column:auth"`
	LastLogin  int64  `json:"last_login" gorm:"column:last_login"`
	LastLogout int64  `json:"last_logout" gorm:"column:last_logout"`
	Created    int64  `json:"created" gorm:"autoCreateTime" gorm:"column:created"`
	Updated    int64  `json:"updated" gorm:"autoUpdateTime:milli" gorm:"column:updated"`
	FlowCount  uint64 `json:"flow_count" gorm:"column:flow_count"`
	Disabled   int    `json:"disabled" gorm:"column:disabled"`
	ConfigId   string `json:"config_id" gorm:"column:config_id"`
}

//
// UserStorage
// @Description:
//
type UserStorage struct {
	RXFlowCounter, TXFlowCounter *traffic.FlowStatisticsFP
	Config                       config.ClientConfig
	Address                      string
	UUID                         string
}

//
// TableName
// @Description: gorm table name
// @receiver UserInfo
// @return string
//
func (UserInfo) TableName() string {
	return "tunn_user"
}

//
// FormatCreate
// @Description:
// @receiver i
//
func (i *UserInfo) FormatCreate() {
	i.Id = ""
	i.Auth = "user"
	i.LastLogin = 0
	i.LastLogout = 0
	i.Created = 0
	i.Updated = 0
	i.FlowCount = 0
	i.ConfigId = ""
}

//
// RemoveSensitive
// @Description:
// @receiver i
//
func (i *UserInfo) RemoveSensitive() {
	i.Password = ""
}

//
// Copy
// @Description:
// @receiver i
// @param src
// @return error
//
func (i *UserInfo) Copy(src UserInfo) error {
	marshal, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, i)
	if err != nil {
		return err
	}
	return nil
}

//
// UserStatus
// @Description: from cache
//
type UserStatus struct {
	Online  bool                `json:"online"`
	Address string              `json:"address"`
	RX      LinkStatus          `json:"rx"`
	TX      LinkStatus          `json:"tx"`
	Config  config.ClientConfig `json:"config"`
	UUID    string              `json:"uuid"`
}

//
// UpdateDiff
// @Description:
// @receiver status
// @param dst
//
func (status *UserStatus) UpdateDiff(dst UserStatus) {
	if status.Online != dst.Online {
		status.Online = dst.Online
	}
	if status.Address != dst.Address {
		status.Address = dst.Address
	}
}
