package administration

import (
	"errors"
	log "github.com/cihub/seelog"
	"gorm.io/gorm"
	"time"
	"tunn-hub/administration/model"
)

//
// userInfoService
// @Description:
//
type userInfoService struct {
	ServerAdmin
}

//
// newUserInfoService
// @Description:
// @param admin
// @return *userInfoService
//
func newUserInfoService(admin *ServerAdmin) *userInfoService {
	return &userInfoService{
		ServerAdmin: *admin,
	}
}

//
// Create
// @Description: 创建用户
// @receiver u
// @param info
// @return error
//
func (u *userInfoService) Create(info *model.UserInfo) error {
	u.EncryptInfo(info)
	if _, err := u.accountExist(info.Account); err != nil {
		return err
	}
	db := u.db.Create(&info)
	if db.Error == nil {
		err := u.Query(info)
		if err != nil {
			return err
		}
	}
	return db.Error
}

//
// accountExist
// @Description: 检查用户是否存在
// @receiver u
// @param account
// @return error
//
func (u *userInfoService) accountExist(account string) (model.UserInfo, error) {
	info := model.UserInfo{}
	db := u.db.Where("account = ?", account).First(&info)
	if db.Error != gorm.ErrRecordNotFound {
		if db.Error == nil {
			return info, errors.New("account already in use")
		}
		return info, db.Error
	}
	return info, nil
}

//
// GetInfoByAccount
// @Description:
// @receiver u
// @param account
// @return model.UserInfo
// @return error
//
func (u *userInfoService) GetInfoByAccount(account string) (model.UserInfo, error) {
	info := model.UserInfo{Account: account}
	u.EncryptInfo(&info)
	db := u.db.Where("account = ?", info.Account).First(&info)
	if db.Error != nil {
		return model.UserInfo{}, db.Error
	}
	return info, nil
}

//
// GetIdByAccount
// @Description: 通过账号查询用户,用于客户端(非admin)登录时注册用户信息使用
// @receiver u
// @param account
// @return string
// @return error
//
func (u *userInfoService) GetIdByAccount(account string) (string, error) {
	info := model.UserInfo{Account: account}
	u.EncryptInfo(&info)
	db := u.db.Where("account = ?", info.Account).First(&info)
	if db.Error != nil {
		return "", db.Error
	}
	return info.Id, nil
}

//
// Delete
// @Description: 删除用户
// @receiver u
// @param info
//
func (u *userInfoService) Delete(info *model.UserInfo) error {
	u.EncryptInfo(info)
	db := u.db.Where("id=?", info.Id).Delete(&info)
	u.DecryptInfo(info)
	return db.Error
}

//
// Save
// @Description: 保存所有数据
// @receiver u
//
func (u *userInfoService) Save(info *model.UserInfo) error {
	u.EncryptInfo(info)
	db := u.db.Where("id=? AND account=?", info.Id, info.Account).Save(&info)
	u.DecryptInfo(info)
	return db.Error
}

//
// Update
// @Description: 保存非零值数据
// @receiver u
// @param info
// @return error
//
func (u *userInfoService) Update(info *model.UserInfo) error {
	u.EncryptInfo(info)
	db := u.db.Model(info).Updates(model.UserInfo{
		Password:   info.Password,
		Email:      info.Email,
		LastLogin:  info.LastLogin,
		LastLogout: info.LastLogout,
		FlowCount:  info.FlowCount,
		TXCount:    info.TXCount,
		Disabled:   info.Disabled,
		ConfigId:   info.ConfigId,
	})
	u.DecryptInfo(info)
	return db.Error
}

//
// UpdateFlowCount
// @Description: 保存流量数据
// @receiver u
// @param id
// @param flow 对应hub向客户端发送的流量
// @param tx 对应hub接收的流量
//
func (u *userInfoService) UpdateFlowCount(id string, flow uint64, tx uint64) {
	//不更新小于100Kb的数据
	if id == "" || (flow <= 100*1024 && tx <= 100*1024) {
		return
	}
	info := model.UserInfo{Id: id}
	db := u.db.Raw("UPDATE tunn_user SET flow_count = ?,tx_count = ?, updated = ? WHERE id = ?",
		gorm.Expr("flow_count + ?", flow), gorm.Expr("tx_count + ?", tx), time.Now().UnixMilli(), info.Id).Scan(&info)
	if db.Error != nil {
		_ = log.Warn("failed to update flow count at [", id, "] : ", db.Error)
		return
	}
}

//
// UpdateByAccount
// @Description: 通过账号保存非零值数据
// @receiver u
// @param info
// @return error
//
func (u *userInfoService) UpdateByAccount(info *model.UserInfo) error {
	u.EncryptInfo(info)
	db := u.db.Model(&info).Where("account=?", info.Account).Updates(model.UserInfo{
		Password:   info.Password,
		Email:      info.Email,
		LastLogin:  info.LastLogin,
		LastLogout: info.LastLogout,
		FlowCount:  info.FlowCount,
		TXCount:    info.TXCount,
		Disabled:   info.Disabled,
		ConfigId:   info.ConfigId,
	})
	u.DecryptInfo(info)
	return db.Error
}

//
// Query
// @Description: 查询
// @receiver u
// @param info
//
func (u *userInfoService) Query(info *model.UserInfo) error {
	u.EncryptInfo(info)
	db := u.db.First(info)
	u.DecryptInfo(info)
	return db.Error
}

//
// CheckPasswordAndAccount
// @Description: 验证用户名密码,供客户端(非webui)登录验证时使用
// @receiver u
// @param info
// @return error
//
func (u *userInfoService) CheckPasswordAndAccount(info model.UserInfo) (cfg model.ClientConfig, err error) {
	cfg = model.ClientConfig{}
	u.EncryptInfo(&info)
	pwd := info.Password
	info.Password = ""
	db := u.db.Raw("SELECT * FROM tunn_user WHERE account=?", info.Account).First(&info)
	if db.Error == nil {
		if info.Password != pwd {
			return cfg, errors.New("password not match")
		} else if info.Disabled == 1 {
			return cfg, errors.New("user disabled")
		} else {
			cfg, err = UserServiceInstance().configService.GetById(info.ConfigId)
			if err != nil {
				return cfg, err
			}
			return
		}
	} else if db.Error == gorm.ErrRecordNotFound {
		return cfg, errors.New("no such user")
	} else {
		return cfg, errors.New("database error")
	}
}

//
// Login
// @Description: 登录(仅admin使用)
// @receiver u
// @param info
// @return error
//
func (u *userInfoService) Login(info *model.UserInfo) error {
	u.EncryptInfo(info)
	pwd := info.Password
	info.Password = ""
	db := u.db.Raw("SELECT * FROM tunn_user WHERE account=?", info.Account).First(&info)
	if db.Error == nil {
		if info.Password != pwd {
			return errors.New("password not match")
		}
		u.DecryptInfo(info)
		return nil
	}
	if db.Error == gorm.ErrRecordNotFound {
		return errors.New("no such user")
	}
	return errors.New("database error")
}

//
// GetByAccount
// @Description: 通过账号获取信息
// @receiver u
// @param account
// @return model.UserInfo
// @return error
//
func (u *userInfoService) GetByAccount(info *model.UserInfo) error {
	u.EncryptInfo(info)
	db := u.db.Raw("SELECT * FROM tunn_user WHERE account=?", info.Account).First(info)
	u.DecryptInfo(info)
	return db.Error
}

//
// List
// @Description: 列表
// @receiver u
// @return []model.UserInfo
//
func (u *userInfoService) List() ([]model.UserInfo, error) {
	infos := make([]model.UserInfo, 0)
	db := u.db.Raw("SELECT * FROM tunn_user").Scan(&infos)
	for i := range infos {
		info := infos[i]
		u.DecryptInfo(&info)
		infos[i] = info
	}
	return infos, db.Error
}

//
// ListFull
// @Description:
// @receiver u
// @return []model.UserFull
// @return error
//
func (u *userInfoService) ListFull() ([]model.UserFull, error) {
	infos := make([]model.UserFull, 0)
	db := u.db.Raw("SELECT tunn_config.content, tunn_user.* FROM tunn_config LEFT JOIN tunn_user ON tunn_config.id=tunn_user.config_id").Scan(&infos)
	for i := range infos {
		info := infos[i]
		u.DecryptFullInfo(&info)
		infos[i] = info
	}
	return infos, db.Error
}

//
// EncryptInfo
// @Description: 加密
// @receiver u
// @param i
// @return model.UserInfo
//
func (u *userInfoService) EncryptInfo(i *model.UserInfo) {
	i.Password = u.crypt.Encrypt(i.Password)
	i.Account = u.crypt.Encrypt(i.Account)
	i.Auth = u.crypt.Encrypt(i.Auth)
	i.Email = u.crypt.Encrypt(i.Email)
}

//
// DecryptInfo
// @Description: 解密
// @receiver u
// @param i
// @return model.UserInfo
//
func (u *userInfoService) DecryptInfo(i *model.UserInfo) {
	i.Password = u.crypt.Decrypt(i.Password)
	i.Account = u.crypt.Decrypt(i.Account)
	i.Auth = u.crypt.Decrypt(i.Auth)
	i.Email = u.crypt.Decrypt(i.Email)
}

//
// DecryptFullInfo
// @Description: 解密
// @receiver u
// @param i
// @return model.UserInfo
//
func (u *userInfoService) DecryptFullInfo(i *model.UserFull) {
	i.Password = u.crypt.Decrypt(i.Password)
	i.Account = u.crypt.Decrypt(i.Account)
	i.Auth = u.crypt.Decrypt(i.Auth)
	i.Email = u.crypt.Decrypt(i.Email)
	storage := model.ClientConfigStorage{
		Id:      i.ConfigId,
		Content: i.Content,
	}
	cfg, err := storage.GetConfig()
	if err != nil {
		return
	}
	i.Config = cfg
}
