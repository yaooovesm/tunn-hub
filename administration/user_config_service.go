package administration

import (
	"errors"
	log "github.com/cihub/seelog"
	"strings"
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
// AvailableExports
// @Description:
// @receiver u
// @return []config.Route
// @return error
//
func (u *userClientConfigService) AvailableExports(account string) ([]model.ImportableRoute, error) {
	routes := make([]model.ImportableRoute, 0)
	serverRoutes := config.Current.Routes
	for i := range serverRoutes {
		if serverRoutes[i].Option == config.RouteOptionExport {
			routes = append(routes, model.ImportableRoute{
				Name:         serverRoutes[i].Name,
				Network:      serverRoutes[i].Network,
				Provider:     "TunnHub",
				Certificated: true,
			})
		}
	}
	list, err := UserServiceInstance().infoService.ListFull()
	//找到自身账号权限
	auth := None
	for i := range list {
		if list[i].Account == account {
			auth = UserAuthLevel(list[i].Auth)
			break
		}
	}
	if err != nil {
		return nil, err
	}
	for i := range list {
		for j := range list[i].Config.Routes {
			if list[i].Config.Routes[j].Option == config.RouteOptionExport {
				if auth == Administrator || visibilityCheck(list[i].Config.Routes[j].Visibility, account) {
					routes = append(routes, model.ImportableRoute{
						Name:         list[i].Config.Routes[j].Name,
						Network:      list[i].Config.Routes[j].Network,
						Provider:     list[i].Account,
						Certificated: list[i].Auth == string(Administrator),
					})
				}
			}
		}
	}
	return routes, nil
}

//
// AllExports
// @Description:
// @receiver u
// @return []model.ImportableRoute
// @return error
//
func (u *userClientConfigService) AllExports() ([]model.ImportableRoute, error) {
	routes := make([]model.ImportableRoute, 0)
	serverRoutes := config.Current.Routes
	for i := range serverRoutes {
		if serverRoutes[i].Option == config.RouteOptionExport {
			routes = append(routes, model.ImportableRoute{
				Name:         serverRoutes[i].Name,
				Network:      serverRoutes[i].Network,
				Provider:     "TunnHub",
				Certificated: true,
			})
		}
	}
	list, err := UserServiceInstance().infoService.ListFull()
	if err != nil {
		return nil, err
	}
	for i := range list {
		for j := range list[i].Config.Routes {
			if list[i].Config.Routes[j].Option == config.RouteOptionExport {
				routes = append(routes, model.ImportableRoute{
					Name:         list[i].Config.Routes[j].Name,
					Network:      list[i].Config.Routes[j].Network,
					Provider:     list[i].Account,
					Certificated: list[i].Auth == string(Administrator),
				})
			}
		}
	}
	return routes, nil
}

//
// visibilityCheck
// @Description:
// @param v
// @param account
// @return bool
//
func visibilityCheck(v, account string) bool {
	if v == "all" {
		return true
	}
	allowed := strings.Split(v, ",")
	for i := range allowed {
		if allowed[i] == account {
			return true
		}
	}
	return false
}

//
// ImportsCheck
// @Description:
// @receiver u
// @param routes
// @return error
//
func (u *userClientConfigService) ImportsCheck(routes []config.Route) error {
	exports, err := u.AllExports()
	if err != nil {
		return err
	}
	for i := range routes {
		r := routes[i]
		//只检查import
		if r.Option != config.RouteOptionImport {
			continue
		}
		hasExport := false
		for j := range exports {
			if exports[j].Network == r.Network {
				hasExport = true
				break
			}
		}
		if !hasExport {
			return errors.New("invalid import " + r.Network + " no export")
		}
	}
	return nil
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
	//添加命名
	exports, err := u.AllExports()
	if err != nil {
		return model.ClientConfig{}, err
	}
	for i := range cfg.Routes {
		for j := range exports {
			if exports[j].Network == cfg.Routes[i].Network {
				cfg.Routes[i].Name = exports[j].Name
			}
		}
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
	err := u.HasDuplicateExport(cfg.Routes, cfg.Id, false)
	if err != nil {
		return model.ClientConfig{}, err
	}
	err = u.ImportsCheck(cfg.Routes)
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
// ResetById
// @Description:
// @receiver u
// @param cfg
// @return error
//
func (u *userClientConfigService) ResetById(cfg model.ClientConfig) error {
	cfg.Routes = []config.Route{}
	cfg.Device = config.Device{}
	storage := cfg.ToStorageModel()
	db := u.db.Where("id=?", cfg.Id).Save(&storage)
	return db.Error
}

//
// Create
// @Description:
// @receiver u
// @param cfg
// @return error
//
func (u *userClientConfigService) Create(cfg model.ClientConfig) (model.ClientConfig, error) {
	err := u.HasDuplicateExport(cfg.Routes, "", false)
	if err != nil {
		return model.ClientConfig{}, err
	}
	err = u.ImportsCheck(cfg.Routes)
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
// HasDuplicateExport
// @Description:
// @receiver u
// @param routes
// @return error
//
func (u *userClientConfigService) HasDuplicateExport(routes []config.Route, ignoreId string, ignoreServer bool) error {
	//自检
	temp := make(map[string]int)
	for i := range routes {
		r := routes[i]
		if _, ok := temp[r.Network]; ok {
			return log.Error("current configuration contains duplicate routes")
		} else {
			temp[r.Network] = 1
		}
	}
	if !ignoreServer {
		//服务端设置的Export
		serverRoutes := config.Current.Routes
		for i := range serverRoutes {
			if serverRoutes[i].Option == config.RouteOptionExport {
				for j := range routes {
					if err := serverRoutes[i].IsDuplicateExport(routes[j]); err != nil {
						return log.Error("export duplicate with server export : ", serverRoutes[i].Network, " , ", routes[j].Network)
					}
				}
			}
		}
	}
	//客户端Export
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
					return log.Error("export duplicate with user export @"+list[i].Id, ": ", list[i].Routes[j].Network, " , ", routes[k].Network)
				}
			}
		}
	}
	return nil
}
