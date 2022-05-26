package administration

import (
	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"time"
	"tunn-hub/administration/model"
	"tunn-hub/config"
)

var userServiceInstance *userService

//
// UserServiceInstance
// @Description:
// @return userService
//
func UserServiceInstance() *userService {
	return userServiceInstance
}

//
// userService
// @Description:
//
type userService struct {
	statusService *userStatusService
	infoService   *userInfoService
	configService *userClientConfigService
}

//
// newUserService
// @Description:
// @param admin
// @return *userService
//
func newUserService(admin *ServerAdmin) *userService {
	serv := &userService{
		statusService: newUserStatusService(admin),
		infoService:   newUserInfoService(admin),
		configService: newUserClientConfigService(admin),
	}
	userServiceInstance = serv
	return serv
}

//
// TunnelLogin
// @Description:
// @receiver serv
// @param account
// @param password
//
func (serv *userService) TunnelLogin(account string, password string) (model.ClientConfig, error) {
	return serv.infoService.CheckPasswordAndAccount(model.UserInfo{Account: account, Password: password})
}

//
// InfoService
// @Description:
// @receiver serv
// @return *userInfoService
//
func (serv *userService) InfoService() *userInfoService {
	return serv.infoService
}

//
// StatusService
// @Description:
// @receiver serv
// @return *userStatusService
//
func (serv *userService) StatusService() *userStatusService {
	return serv.statusService
}

//
// Login
// @Description:
// @receiver serv
//
func (serv *userService) Login(ctx *gin.Context, info *model.UserInfo) (string, error) {
	err := serv.infoService.Login(info)
	if err != nil {
		return "", err
	}
	remote, _ := ctx.RemoteIP()
	return TokenServiceInstance().createByInfo(remote, *info), nil
}

//
// SetOnline
// @Description:
// @receiver serv
// @param account
// @param address
//
func (serv *userService) SetOnline(account string, address string) {
	if account == "" {
		return
	}
	serv.SetStatusByAccount(account, model.UserStatus{
		Online:  true,
		Address: address,
	})
	info := model.UserInfo{
		Account:   account,
		LastLogin: time.Now().UnixMilli(),
	}
	err := serv.infoService.UpdateByAccount(&info)
	if err != nil {
		_ = log.Warn("error during update : ", err)
		return
	}
}

//
// SetOffline
// @Description:
// @receiver serv
// @param account
//
func (serv *userService) SetOffline(account string) {
	if account == "" {
		return
	}
	serv.SetStatusByAccount(account, model.UserStatus{
		Online:  false,
		Address: "",
		Config:  config.Config{},
		RX:      model.LinkStatus{},
		TX:      model.LinkStatus{},
	})
	//删除储存
	serv.statusService.DeleteStorage(account)
	info := model.UserInfo{
		Account:    account,
		LastLogout: time.Now().UnixMilli(),
	}
	err := serv.infoService.UpdateByAccount(&info)
	if err != nil {
		_ = log.Warn("error during update : ", err)
		return
	}
}

//
// SetStatusByAccount
// @Description: 更新所有与缓存中不同的数据
// @receiver serv
//
func (serv *userService) SetStatusByAccount(account string, status model.UserStatus) {
	id, err := serv.infoService.GetIdByAccount(account)
	if err != nil {
		return
	}
	serv.statusService.UpdateStatus(id, status)
}

//
// GetUserInfoById
// @Description:
// @receiver serv
// @param account
// @return model.UserInfo
//
func (serv userService) GetUserInfoById(id string) (model.UserInfo, error) {
	info := model.UserInfo{Id: id}
	err := serv.infoService.Query(&info)
	return info, err
}

//
// GetUserInfoByAccount
// @Description:
// @receiver serv
// @param account
// @return model.UserInfo
//
func (serv userService) GetUserInfoByAccount(account string) (model.UserInfo, error) {
	info := model.UserInfo{Account: account}
	err := serv.infoService.Query(&info)
	return info, err
}

//
// GetUserByAccount
// @Description:
// @receiver serv
// @param account
//
func (serv userService) GetUserByAccount(account string) (model.User, error) {
	info := model.UserInfo{Account: account}
	err := serv.infoService.GetByAccount(&info)
	if err != nil {
		return model.User{}, err
	}
	return model.User{
		UserInfo:   info,
		UserStatus: serv.statusService.GetStatus(info.Id),
	}, nil
}

//
// GetUserById
// @Description:
// @receiver serv
// @param id
//
func (serv userService) GetUserById(id string) (model.User, error) {
	info := model.UserInfo{Id: id}
	err := serv.infoService.Query(&info)
	if err != nil {
		return model.User{}, err
	}
	return model.User{
		UserInfo:   info,
		UserStatus: serv.statusService.GetStatus(info.Id),
	}, nil
}

//
// GetGeneral
// @Description:
// @receiver serv
// @return map[string]int
// @return error
//
func (serv userService) GetGeneral() (map[string]int, error) {
	users, err := serv.ListUsers()
	if err != nil {
		return nil, err
	}
	total := 0
	online := 0
	disabled := 0
	for i := range users {
		user := users[i]
		if user.Id == "" {
			continue
		}
		total++
		if user.Online {
			online++
		}
		if user.Disabled == 1 {
			disabled++
		}
	}
	return map[string]int{
		"total":    total,
		"online":   online,
		"disabled": disabled,
	}, nil
}

//
// ListUsers
// @Description:
// @receiver serv
// @return []model.User
//
func (serv userService) ListUsers() (list []model.User, err error) {
	infos, err := serv.infoService.List()
	if err != nil {
		return nil, err
	}
	for i := range infos {
		list = append(list, model.User{
			UserInfo:   infos[i],
			UserStatus: serv.statusService.GetStatus(infos[i].Id),
		})
	}
	return
}

//
// ListUserInfo
// @Description:
// @receiver serv
// @return []model.User
//
func (serv userService) ListUserInfo() (list []model.UserInfo, err error) {
	return serv.infoService.List()
}

//
// CreateUser
// @Description:
// @receiver serv
// @param info
//
func (serv userService) CreateUser(info *model.UserInfo) error {
	id := UUID()
	info.Id = id
	cfg, err := serv.configService.Create(model.ClientConfig{
		Routes: []config.Route{},
		Device: config.Device{},
	})
	if err != nil {
		return err
	}
	info.ConfigId = cfg.Id
	err = serv.infoService.Create(info)
	if err != nil {
		return err
	}
	return nil
}

//
// UpdateUser
// @Description:
// @receiver serv
// @param info
//
func (serv userService) UpdateUser(info *model.UserInfo) error {
	origin := model.UserInfo{Id: info.Id}
	err := serv.infoService.Query(&origin)
	if err != nil {
		return err
	}
	origin.Email = info.Email
	origin.ConfigId = info.ConfigId
	if info.Password != "" {
		origin.Password = info.Password
	}
	err = serv.infoService.Save(&origin)
	if err != nil {
		return err
	}
	_ = info.Copy(origin)
	return nil
}

//
// SetUserDisable
// @Description:
// @receiver serv
// @param info
// @return error
//
func (serv userService) SetUserDisable(info *model.UserInfo) error {
	origin := model.UserInfo{Id: info.Id}
	err := serv.infoService.Query(&origin)
	if err != nil {
		return err
	}
	origin.Disabled = info.Disabled
	err = serv.infoService.Save(&origin)
	if err != nil {
		return err
	}
	_ = info.Copy(origin)
	return nil
}

//
// ResetFlowCounter
// @Description:
// @receiver serv
// @param info
// @return error
//
func (serv userService) ResetFlowCounter(info *model.UserInfo) error {
	origin := model.UserInfo{Id: info.Id}
	err := serv.infoService.Query(&origin)
	if err != nil {
		return err
	}
	origin.FlowCount = 0
	err = serv.infoService.Save(&origin)
	if err != nil {
		return err
	}
	_ = info.Copy(origin)
	return nil
}

//
// DeleteUser
// @Description:
// @receiver serv
//
func (serv userService) DeleteUser(info *model.UserInfo) error {
	err := serv.infoService.Query(info)
	if err != nil {
		return err
	}
	err = serv.infoService.Delete(info)
	if err != nil {
		return err
	}
	return nil
}

//
// GrantAdmin
// @Description:
// @receiver serv
//
func (serv userService) GrantAdmin(info *model.UserInfo) error {
	origin := model.UserInfo{Id: info.Id}
	err := serv.infoService.Query(&origin)
	if err != nil {
		return err
	}
	origin.Auth = string(Administrator)
	err = serv.infoService.Save(&origin)
	if err != nil {
		return err
	}
	_ = info.Copy(origin)
	return nil
}

//
// GenerateUserConfig
// @Description:
// @receiver serv
// @param id
// @return config.PushedConfig
// @return error
//
func (serv userService) GenerateUserConfig(id string) (config.PushedConfig, error) {
	info := model.UserInfo{Id: id}
	err := serv.infoService.Query(&info)
	if err != nil {
		return config.PushedConfig{}, err
	}
	clientConfig, err := serv.configService.GetById(info.ConfigId)
	if err != nil {
		return config.PushedConfig{}, err
	}
	return clientConfig.ToPushModel(), nil
}

//
// CommitAllFlowCount
// @Description:
// @receiver serv
//
func (serv *userService) CommitAllFlowCount() {
	serv.statusService.CommitAllFlowCounts()
}

//
// AvailableExports
// @Description:
// @receiver serv
// @return []config.Route
// @return error
//
func (serv *userService) AvailableExports() ([]config.Route, error) {
	return serv.configService.AvailableExports()
}
