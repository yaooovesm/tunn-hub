package administration

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net"
	"sync"
	"time"
	"tunn-hub/administration/model"
)

var tokenServiceInstance *TokenService

//
// TokenServiceInstance
// @Description:
// @return TokenService
//
func TokenServiceInstance() *TokenService {
	return tokenServiceInstance
}

const (
	DefaultTokenTimeout       = 30 * 60 * 1000
	DefaultTokenCleanInterval = 60 * 1000
)

var (
	ErrNotLogin               = errors.New("user not login")
	ErrTokenExpired           = errors.New("token expired")
	ErrRemoteHostNotPermitted = errors.New("remote host not permitted")
	ErrNoPermissions          = errors.New("no permissions")
)

type UserAuthLevel string

const (
	Administrator UserAuthLevel = "admin"
	User          UserAuthLevel = "user"
	None          UserAuthLevel = ""
)

var AuthLevelMap = map[UserAuthLevel]int{
	Administrator: 100,
	User:          1,
	None:          0,
}

//
// Compare
// @Description:
// @receiver l
// @param l2
// @return high
//
func (l UserAuthLevel) Compare(l2 UserAuthLevel) (high UserAuthLevel) {
	if v1, ok := AuthLevelMap[l]; !ok {
		return l2
	} else {
		if v2, ok := AuthLevelMap[l2]; !ok {
			return l
		} else {
			if v1 > v2 {
				return l
			}
			return l2
		}
	}
}

//
// TokenService
// @Description:
//
type TokenService struct {
	crypt Crypt
	cache sync.Map
	exp   sync.Map
	ip    sync.Map
	level sync.Map
}

//
// newTokenServer
// @Description:
// @param crypt
//
func newTokenServer(crypt Crypt) *TokenService {
	server := &TokenService{
		crypt: crypt,
		cache: sync.Map{},
		exp:   sync.Map{},
		level: sync.Map{},
		ip:    sync.Map{},
	}
	server.autoClean()
	tokenServiceInstance = server
	return server
}

//
// autoClean
// @Description:
// @receiver c
//
func (t *TokenService) autoClean() {
	go func() {
		interval := time.Millisecond * time.Duration(DefaultTokenCleanInterval)
		for {
			current := time.Now().UnixMilli()
			t.exp.Range(func(key, value interface{}) bool {
				if value.(int64) < current {
					t.RemoveByCode(key.(string))
				}
				return true
			})
			time.Sleep(interval)
		}
	}()
}

//
// CheckHeader
// @Description:
// @receiver t
// @param ctx
// @param level
// @return error
//
func (t *TokenService) CheckHeader(ctx *gin.Context, level UserAuthLevel) error {
	return t.CheckToken(ctx, level)
}

//
// CreateAndCheck
// @Description:
// @receiver t
// @param account
//
func (t *TokenService) CreateAndCheck(remote net.IP, account string, password string) (string, error) {
	info := model.UserInfo{Account: account, Password: password}
	err := UserServiceInstance().infoService.Login(&info)
	if err != nil {
		return "", err
	}
	return t.createByInfo(remote, info), nil
}

//
// RemoveByCode
// @Description:
// @receiver t
// @param code
//
func (t *TokenService) RemoveByCode(code string) {
	t.cache.Delete(code)
	t.cache.Delete(code + "_id")
	t.exp.Delete(code)
	t.level.Delete(code)
	t.ip.Delete(code)
}

//
// createByInfo
// @Description:
// @receiver t
// @param info
//
func (t *TokenService) createByInfo(remote net.IP, info model.UserInfo) string {
	token := model.Token{
		Account: info.Account,
		UUID:    UUID(),
		Id:      info.Id,
	}
	code := t.encodeToken(token)
	t.cache.Store(code, info.Account)
	t.cache.Store(code+"_id", info.Id)
	t.exp.Store(code, time.Now().UnixMilli()+DefaultTokenTimeout)
	t.level.Store(code, info.Auth)
	t.ip.Store(code, remote)
	return code
}

//
// encodeToken
// @Description:
// @receiver t
// @param token
// @return string
//
func (t *TokenService) encodeToken(token model.Token) string {
	marshal, _ := json.Marshal(token)
	return base64.StdEncoding.EncodeToString([]byte(t.crypt.Encrypt(string(marshal))))
}

//
// decodeToken
// @Description:
// @receiver t
// @param code
// @return model.Token
//
func (t *TokenService) decodeToken(code string) model.Token {
	b, _ := base64.StdEncoding.DecodeString(code)
	decryptStr := t.crypt.Decrypt(string(b))
	token := model.Token{}
	_ = json.Unmarshal([]byte(decryptStr), &token)
	return token
}

//
// IsExist
// @Description:
// @receiver t
// @param account
// @param code
// @return error
//
func (t *TokenService) IsExist(account string, code string) error {
	if !t.hasToken(code) {
		return ErrNotLogin
	}
	if t.isExpired(code) {
		return ErrTokenExpired
	}
	if value, ok := t.cache.Load(code); ok && value != nil {
		if value.(string) != account {
			return errors.New("operation not allowed")
		}
	}
	return nil
}

//
// CheckToken
// @Description:
// @receiver t
// @param account
//
func (t *TokenService) CheckToken(ctx *gin.Context, level UserAuthLevel) error {
	code := ctx.GetHeader("token")
	ip, _ := ctx.RemoteIP()
	return t.CheckTokenCode(code, level, ip)
}

//
// CheckTokenCode
// @Description:
// @receiver t
// @param code
// @param level
// @param remote
// @return error
//
func (t *TokenService) CheckTokenCode(code string, level UserAuthLevel, remote net.IP) error {
	if !t.hasToken(code) {
		return ErrNotLogin
	}
	if t.isExpired(code) {
		return ErrTokenExpired
	}
	if !t.isPermittedIp(code, remote) {
		return ErrRemoteHostNotPermitted
	}
	if !t.hasPermissions(code, level) {
		return ErrNoPermissions
	}
	//续期
	t.exp.Store(code, time.Now().UnixMilli()+DefaultTokenTimeout)
	return nil
}

//
// CheckSelfByAccount
// @Description:
// @receiver t
//
func (t *TokenService) CheckSelfByAccount(ctx *gin.Context, account string) error {
	code := ctx.GetHeader("token")
	if level, ok := t.level.Load(code); ok && level != nil {
		if UserAuthLevel(level.(string)) == Administrator {
			return nil
		}
	}
	if value, ok := t.cache.Load(code); ok && value != nil {
		if value.(string) != account {
			return errors.New("operation not allowed")
		}
	}
	return nil
}

//
// CheckSelfById
// @Description:
// @receiver t
//
func (t *TokenService) CheckSelfById(ctx *gin.Context, id string) error {
	code := ctx.GetHeader("token")
	if level, ok := t.level.Load(code); ok && level != nil {
		if UserAuthLevel(level.(string)) == Administrator {
			return nil
		}
	}
	if value, ok := t.cache.Load(code + "_id"); ok && value != nil {
		if value.(string) != id {
			return errors.New("operation not allowed")
		}
	}
	return nil
}

//
// hasPermissions
// @Description:
// @receiver t
// @param code
// @param level
//
func (t *TokenService) hasPermissions(code string, level UserAuthLevel) bool {
	has := false
	if lev, ok := AuthLevelMap[level]; ok {
		t.level.Range(func(key, value interface{}) bool {
			if key.(string) == code {
				cur := AuthLevelMap[UserAuthLevel(value.(string))]
				if cur >= lev {
					has = true
				}
				return false
			}
			return true
		})
	}
	return has
}

//
// hasToken
// @Description:
// @receiver t
// @param code
// @return bool
//
func (t *TokenService) hasToken(code string) bool {
	has := false
	t.cache.Range(func(key, value interface{}) bool {
		if key.(string) == code {
			has = true
			return false
		}
		return true
	})
	return has
}

//
// isExpired
// @Description:
// @receiver t
// @param code
// @return bool
//
func (t *TokenService) isExpired(code string) bool {
	expired := true
	current := time.Now().UnixMilli()
	t.exp.Range(func(key, value interface{}) bool {
		if key.(string) == code {
			if value.(int64) > current {
				//通过则延时
				expired = false
				t.exp.Store(key, time.Now().UnixMilli()+DefaultTokenTimeout)
			}
			return false
		}
		return true
	})
	return expired
}

//
// isPermittedIp
// @Description:
// @receiver t
// @param ip
//
func (t *TokenService) isPermittedIp(code string, ip net.IP) bool {
	if value, ok := t.ip.Load(code); ok {
		if value == nil && ip == nil ||
			value.(net.IP).String() == ip.String() {
			return true
		}
		return false
	} else {
		return false
	}
}
