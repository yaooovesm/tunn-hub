package authentication

import (
	"encoding/hex"
	log "github.com/cihub/seelog"
	"github.com/gofrs/uuid"
	"tunn-hub/common/config"
)

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
	config.Current.DataProcess.Key = key
	log.Info("[cipher:", len(key), "] update key")
	log.Info("convert to hex -> ", hex.EncodeToString(key))
}
