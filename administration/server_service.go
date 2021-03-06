package administration

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"
	"tunn-hub/administration/model"
	"tunn-hub/config"
	"tunn-hub/monitor"
	"tunn-hub/networking"
	"tunn-hub/traffic"
	"tunn-hub/transmitter"
)

var serverServiceInstance *serverService

//
// ServerServiceInstance
// @Description:
// @return *serverService
//
func ServerServiceInstance() *serverService {
	return serverServiceInstance
}

//
// newServerService
// @Description:
// @return *serverService
//
func newServerService() *serverService {
	serverServiceInstance = &serverService{
		monitorService:  newMonitorService(4000),
		recorderService: newRecorderService(),
	}
	return serverServiceInstance
}

//
// serverService
// @Description:
//
type serverService struct {
	recorderService    *recorderService
	monitorService     *monitorService
	rxFlowCounter      *traffic.FlowStatisticsFP
	txFlowCounter      *traffic.FlowStatisticsFP
	ctx                context.Context
	cancel             context.CancelFunc
	kickFunc           func(uuid string) error
	reconnectFunc      func(uuid string) error
	searchFunc         func(uuid string) (cfg config.ClientConfig, err error)
	transmitterVersion transmitter.Version
	wsKey              string
	ippool             *networking.IPAddressPool
}

//
// SetupServer
// @Description:
// @receiver serv
// @param rx
// @param tx
//
func (serv *serverService) SetupServer(rx, tx *traffic.FlowStatisticsFP,
	context context.Context, cancel context.CancelFunc) {
	serv.rxFlowCounter = rx
	serv.txFlowCounter = tx
	serv.ctx = context
	serv.cancel = cancel
}

//
// SetupAuthServer
// @Description:
// @receiver serv
// @param kickFunc
//
func (serv *serverService) SetupAuthServer(
	kickFunc func(uuid string) error,
	reconnectFunc func(uuid string) error,
	searchFunc func(uuid string) (cfg config.ClientConfig, err error),
	version transmitter.Version,
	wsKey string,
	ippool *networking.IPAddressPool) {
	serv.kickFunc = kickFunc
	serv.reconnectFunc = reconnectFunc
	serv.searchFunc = searchFunc
	serv.transmitterVersion = version
	serv.wsKey = wsKey
	serv.ippool = ippool
}

//
// GetFlowStatus
// @Description:
// @receiver serv
// @return model.ServerStatus
//
func (serv *serverService) GetFlowStatus() model.ServerStatus {
	rx := model.LinkStatus{}
	rx.ReadFromFP(serv.rxFlowCounter)
	tx := model.LinkStatus{}
	tx.ReadFromFP(serv.txFlowCounter)
	return model.ServerStatus{
		RX: rx,
		TX: tx,
	}
}

//
// GetServerConfigs
// @Description:
// @receiver serv
// @return config.Config
//
func (serv *serverService) GetServerConfigs() config.Config {
	return config.Current
}

//
// KickById
// @Description:
// @receiver serv
// @param id
// @return error
//
func (serv *serverService) KickById(id string) error {
	status := UserServiceInstance().statusService.GetStatus(id)
	if status.Online && status.UUID != "" {
		return serv.kickFunc(status.UUID)
	}
	return errors.New("user not online")
}

//
// ReconnectById
// @Description:
// @receiver serv
// @param id
// @return error
//
func (serv *serverService) ReconnectById(id string) error {
	status := UserServiceInstance().statusService.GetStatus(id)
	if status.Online && status.UUID != "" {
		return serv.reconnectFunc(status.UUID)
	}
	return errors.New("user not online")
}

//
// GetIPPoolAllocInfo
// @Description:
// @receiver serv
// @return map[string]networking.IPAllocInfo
//
func (serv *serverService) GetIPPoolAllocInfo() []interface{} {
	if serv.ippool == nil {
		return nil
	}
	info := serv.ippool.Info()
	var result = make([]interface{}, 0)
	for i := range info {
		cfg, err := serv.searchFunc(info[i].UUID)
		if err != nil {
			continue
		}
		result = append(result, map[string]interface{}{
			"account": cfg.User.Account,
			"info":    info[i],
		})
	}
	return result
}

//
// GetIPPoolGeneral
// @Description:
// @receiver serv
// @return map[string]interface{}
//
func (serv *serverService) GetIPPoolGeneral() map[string]interface{} {
	if serv.ippool == nil {
		return nil
	}
	return serv.ippool.General()
}

//
// GetRecentHourTrafficData
// @Description:
// @receiver serv
// @param recent
// @return []monitor.HubTrafficStamp
//
func (serv *serverService) GetRecentHourTrafficData() []monitor.HubTrafficStamp {
	return serv.recorderService.RecentHour()
}

//
// GetRecentDayTrafficData
// @Description:
// @receiver serv
//
func (serv *serverService) GetRecentDayTrafficData() []monitor.HubTrafficStamp {
	return serv.recorderService.RecentDay()
}

//
// readTrafficHistory
// @Description:
// @param path
// @param num
// @param now
// @param duration
// @return []monitor.HubTrafficStamp
//
func readTrafficHistory(path string, num int, now, duration int64) []monitor.HubTrafficStamp {
	data := make([]monitor.HubTrafficStamp, 0)
	read, err := ioutil.ReadFile(path + "traffic" + strconv.Itoa(num))
	if err != nil {
		return make([]monitor.HubTrafficStamp, 0)
	}
	err = json.Unmarshal(read, &data)
	if err != nil {
		return make([]monitor.HubTrafficStamp, 0)
	}
	//????????????24???????????????
	last := len(data) - 1
	if last < 0 || now-data[last].Timestamp > duration {
		return make([]monitor.HubTrafficStamp, 0)
	}
	return data
}
