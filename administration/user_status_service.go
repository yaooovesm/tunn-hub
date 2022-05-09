package administration

import (
	"time"
	"tunn-hub/administration/model"
	"tunn-hub/config"
)

//
// userStatusService
// @Description:
//
type userStatusService struct {
	ServerAdmin
	cache       *Cache
	cachePrefix string
	users       map[string]*model.UserStorage
}

//
// newUserStatusService
// @Description:
// @param admin
// @return *userStatusService
//
func newUserStatusService(admin *ServerAdmin) *userStatusService {
	serv := &userStatusService{
		ServerAdmin: *admin,
		cache:       NewCache(0),
		users:       make(map[string]*model.UserStorage),
	}
	serv.autoCommit()
	return serv
}

//
// autoCommit
// @Description:
// @receiver u
//
func (u *userStatusService) autoCommit() {
	go func() {
		for {
			time.Sleep(time.Minute * 1)
			for id := range u.users {
				if u.users[id] != nil && u.users[id].TXFlowCounter != nil {
					UserServiceInstance().infoService.UpdateFlowCount(id, u.users[id].TXFlowCounter.Commit())
				}
			}
		}
	}()
}

//
// prefix
// @Description:
// @receiver u
// @param id
// @return string
//
func (u *userStatusService) prefix(id string) string {
	return u.cachePrefix + id
}

//
// UpdateStatus
// @Description:
// @receiver u
// @param id
//
func (u *userStatusService) UpdateStatus(id string, status model.UserStatus) {
	cached := u.GetStatus(u.prefix(id))
	cached.UpdateDiff(status)
	u.cache.Set(u.prefix(id), cached)
}

//
// GetStatus
// @Description:
// @receiver u
// @param id
// @return interface{}
//
func (u *userStatusService) GetStatus(id string) model.UserStatus {
	cached := u.cache.Get(u.prefix(id))
	if cached == nil {
		return model.UserStatus{
			Online:  false,
			Address: "",
			RX:      model.LinkStatus{},
			TX:      model.LinkStatus{},
			Config:  config.Config{},
			UUID:    "",
		}
	}
	//收集rx数据
	rx := model.LinkStatus{
		Flow:        0,
		FlowSpeed:   0,
		Packet:      0,
		PacketSpeed: 0,
	}
	//收集tx数据
	tx := model.LinkStatus{
		Flow:        0,
		FlowSpeed:   0,
		Packet:      0,
		PacketSpeed: 0,
	}
	cfg := config.Config{}
	uuid := ""
	if storage, ok := u.users[id]; ok && storage != nil {
		rx.ReadFromFP(storage.RXFlowCounter)
		tx.ReadFromFP(storage.TXFlowCounter)
		cfg = storage.Config
		uuid = storage.UUID
	}

	status := cached.(model.UserStatus)
	//合入数据
	status.RX = rx
	status.TX = tx
	status.Config = cfg
	status.UUID = uuid
	return status
}

//
// RegisterStorage
// @Description:
// @receiver u
// @param account
// @param uss
//
func (u *userStatusService) RegisterStorage(account string, uss *model.UserStorage) {
	id, err := UserServiceInstance().infoService.GetIdByAccount(account)
	if err != nil {
		return
	}
	u.users[id] = uss
}

//
// DeleteStorage
// @Description:
// @receiver u
// @param account
//
func (u *userStatusService) DeleteStorage(account string) {
	id, err := UserServiceInstance().infoService.GetIdByAccount(account)
	if err != nil {
		return
	}
	if u.users[id] != nil && u.users[id].TXFlowCounter != nil {
		UserServiceInstance().infoService.UpdateFlowCount(id, u.users[id].TXFlowCounter.Commit())
	}
	delete(u.users, id)
}
