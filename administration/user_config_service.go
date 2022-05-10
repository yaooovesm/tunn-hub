package administration

import (
	"errors"
	"tunn-hub/administration/model"
	"tunn-hub/config"
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
	if id == "" {
		return model.ClientConfig{}, errors.New("config not found")
	}
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
		cfg, err := list[i].GetConfig()
		if err != nil {
			continue
		}
		cfgs = append(cfgs, cfg)
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
	err := u.HasDuplicateExport(cfg.Routes, cfg.Id)
	if err != nil {
		return model.ClientConfig{}, err
	}
	storage := cfg.ToStorageModel()
	db := u.db.Where("id=?", cfg.Id).Save(&storage)
	if db.Error != nil {
		return model.ClientConfig{}, db.Error
	}
	cfg, err = storage.GetConfig()
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
	err := u.HasDuplicateExport(cfg.Routes, "")
	if err != nil {
		return model.ClientConfig{}, err
	}
	cfg.Id = UUID()
	storage := cfg.ToStorageModel()
	db := u.db.Create(&storage)
	if db.Error != nil {
		return model.ClientConfig{}, db.Error
	}
	err = cfg.Decode(storage.Content)
	if err != nil {
		return model.ClientConfig{}, err
	}
	return cfg, nil
}

//
// IsDuplicateExport
// @Description:
// @receiver u
// @param route
// @return error
//
func (u *userClientConfigService) IsDuplicateExport(route config.Route) error {
	//不在任何一个已暴露的子网内
	//不等于任何一个子网
	list, err := u.List()
	if err != nil {
		return err
	}
	for i := range list {
		for j := range list[i].Routes {
			if err := list[i].Routes[j].IsDuplicateExport(route); err != nil {
				return err
			}
		}
	}
	return nil
}

//
// HasDuplicateExport
// @Description:
// @receiver u
// @param routes
// @return error
//
func (u *userClientConfigService) HasDuplicateExport(routes []config.Route, ignoreId string) error {
	list, err := u.List()
	if err != nil {
		return err
	}
	for i := range list {
		if ignoreId != "" && list[i].Id == ignoreId {
			continue
		}
		for j := range list[i].Routes {
			for k := range routes {
				if err := list[i].Routes[j].IsDuplicateExport(routes[k]); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
