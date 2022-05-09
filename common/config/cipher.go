package config

import (
	log "github.com/cihub/seelog"
	"github.com/gofrs/uuid"
)

type CipherType string

const (
	CipherTypeAES256   CipherType = "AES256"
	CipherTypeAES192   CipherType = "AES192"
	CipherTypeAES128   CipherType = "AES128"
	CipherTypeXOR      CipherType = "XOR"
	CipherTypeSM4      CipherType = "SM4"
	CipherTypeTEA      CipherType = "TEA"
	CipherTypeXTEA     CipherType = "XTEA"
	CipherTypeSalsa20  CipherType = "Salsa20"
	CipherTypeBlowfish CipherType = "Blowfish"
	CipherTypeNone     CipherType = ""
)

//
// DataProcess
// @Description:
//
type DataProcess struct {
	CipherType CipherType `json:"encrypt"`
	Key        []byte     `json:"key"`
}

//
// GenerateCipherKey
// @Description:
//
func GenerateCipherKey() {
	var key []byte
	p1, _ := uuid.NewV4()
	p2, _ := uuid.NewV4()
	key = append(key, p1.Bytes()...)
	key = append(key, p2.Bytes()...)
	Current.DataProcess.Key = key
	log.Info("[cipher:", len(key), "] update key")
}
