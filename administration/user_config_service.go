package administration

import (
	"tunn-hub/administration/model"
)

//
// userClientConfigService
// @Description:
//
type userClientConfigService struct {
	ServerAdmin
}

//
// newUserClientConfigService
// @Description:
// @param admin
// @return *userClientConfigService
//
func newUserClientConfigService(admin *ServerAdmin) *userClientConfigService {
	return &userClientConfigService{
		ServerAdmin: *admin,
	}
}

//
// GetById
// @Description:
// @receiver u
// @param id
// @return model.ClientConfig
//
func (u *userClientConfigService) GetById(id string) (cfg model.ClientConfig, err error) {
	storage := model.ClientConfigStorage{}
	db := u.db.Raw("SELECT * From tunn_config WHERE id=?", id).First(&storage)
	if db.Error != nil {
		return model.ClientConfig{}, db.Error
	}
	cfg, err = storage.GetConfig()
	if err != nil {
		return model.ClientConfig{}, err
	}
	return
}

//
// List
// @Description:
// @receiver u
// @return []model.ClientConfig
// @return error
//
func (u *userClientConfigService) List() ([]model.ClientConfig, error) {
	list := make([]model.ClientConfigStorage, 0)
	db := u.db.Raw("SELECT * FROM tunn_config").Scan(&list)
	if db.Error != nil {
		return []model.ClientConfig{}, db.Error
	}
	cfgs := make([]model.ClientConfig, 0)
	for i := range list {
		config, err := list[i].GetConfig()
		if err != nil {
			continue
		}
		cfgs = append(cfgs, config)
	}
	return cfgs, nil
}

//
// DeleteById
// @Description:
// @receiver u
// @param id
// @return error
//
func (u *userClientConfigService) DeleteById(id string) error {
	storage := model.ClientConfigStorage{Id: id}
	db := u.db.Where("id=?", id).Delete(&storage)
	return db.Error
}

//
// UpdateById
// @Description:
// @receiver u
// @param id
// @param cfg
// @return error
//
func (u *userClientConfigService) UpdateById(cfg model.ClientConfig) (model.ClientConfig, error) {
	storage := cfg.ToStorageModel()
	db := u.db.Where("id=?", cfg.Id).Save(&storage)
	if db.Error != nil {
		return model.ClientConfig{}, db.Error
	}
	cfg, err := storage.GetConfig()
	if err != nil {
		return model.ClientConfig{}, err
	}
	return cfg, nil
}

//
// Create
// @Description:
// @receiver u
// @param cfg
// @return error
//
func (u *userClientConfigService) Create(cfg model.ClientConfig) (model.ClientConfig, error) {
	cfg.Id = UUID()
	storage := cfg.ToStorageModel()
	db := u.db.Create(&storage)
	if db.Error != nil {
		return model.ClientConfig{}, db.Error
	}
	err := cfg.Decode(storage.Content)
	if err != nil {
		return model.ClientConfig{}, err
	}
	return cfg, nil
}
