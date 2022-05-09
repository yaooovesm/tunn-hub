package authentication

import (
	"encoding/json"
	"tunn-hub/common/config"
)

type OptionCode int

const (
	OptionCodeLogin  OptionCode = 0
	OptionCodeLogout OptionCode = 1
)

//
// Info
// @Description:
//
type Info struct {
	AuthInfo AuthInfo `json:"info"`
	Address  string   `json:"address"`
}

//
// AuthInfo
// @Description:
//
type AuthInfo struct {
	Config config.Config `json:"config"`
	UUID   string        `json:"UUID"`
	Option OptionCode
	Bytes  []byte
}

//
// Encode
// @Description:
// @receiver r
//
func (i *AuthInfo) Encode() []byte {
	info := i
	marshal, _ := json.Marshal(info)
	i.Bytes = marshal
	return i.Bytes
}

//
// Decode
// @Description:
// @receiver r
//
func (i *AuthInfo) Decode() error {
	info := &AuthInfo{}
	err := json.Unmarshal(i.Bytes, info)
	if err != nil {
		return err
	}
	i.Config = info.Config
	i.UUID = info.UUID
	i.Option = info.Option
	return nil
}
