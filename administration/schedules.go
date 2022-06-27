package administration

import (
	"encoding/json"
	log "github.com/cihub/seelog"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"tunn-hub/administration/model"
	"tunn-hub/config"
	"tunn-hub/monitor"
)

//
// setupSchedules
// @Description:
// @param service
//
func setupSchedules(service *scheduleService) {
	//每小时保存一次流量记录,总共保留24小时内记录
	err := service.Register("@hourly", "hourly_traffic_recorder", func() {
		//确定要保存的文件名
		//找到最新的序号，最新序号+1即为保存的序号
		path := "./record/hourly/"
		if config.Current.Schedule.HourlyNetworkTrafficDump != "" {
			path = config.Current.Schedule.HourlyNetworkTrafficDump
		}
		s, err := os.Stat(path)
		if os.IsNotExist(err) || !s.IsDir() {
			_ = os.MkdirAll(path, 0644)
		}
		var latest = int64(0)
		var latestNum = 0
		var has = false
		err = filepath.Walk(path,
			func(p string, f os.FileInfo, err error) error {
				if f != nil && !f.IsDir() {
					mo := f.ModTime().UnixMilli()
					num, err := strconv.Atoi(strings.ReplaceAll(filepath.Base(p), "traffic", ""))
					if err != nil {
						return err
					}
					if mo > latest {
						latest = mo
						latestNum = num
						has = true
					}
				}
				return nil
			})
		if err != nil {
			_ = log.Warn("hourly record failed : ", err)
			return
		}
		//确认序号写入
		if has {
			if latestNum == 23 {
				latestNum = 0
			} else {
				latestNum += 1
			}
		}
		err = ServerServiceInstance().recorderService.TrafficRecorder.DumpAndReset(path + "traffic" + strconv.Itoa(latestNum))
		if err != nil {
			_ = log.Warn("hourly record failed : ", err)
			return
		}
	}, true)
	//每分钟记录一次流量记录
	err = service.Register("*/1 * * * *", "traffic_recorder", func() {
		ServerServiceInstance().recorderService.TrafficRecorder.Push(
			monitor.CreateHubTrafficStampWithDuration(
				ServerServiceInstance().rxFlowCounter,
				ServerServiceInstance().txFlowCounter,
				60000, //记录时长 60000ms
			),
		)
	}, true)
	if err != nil {
		_ = log.Warn("schedule register failed : ", err)
	}
	//网络导入检查
	err = service.Register("*/1 * * * *", "user_imports_check_1m", func() {
		checkUserImports()
	}, true)
	if err != nil {
		_ = log.Warn("schedule register failed : ", err)
	}
	//日度流量记录
	if config.Current.Schedule.DailyFlowRecord {
		err := service.Register("@daily", "daily_flow_recorder", func() {
			storagePath := "./record/daily/"
			s, err := os.Stat(storagePath)
			if os.IsNotExist(err) || !s.IsDir() {
				_ = os.MkdirAll(storagePath, 0644)
			}
			err = dumpFlowData(storagePath+"daily_"+time.Now().Format("2006_01_02")+".record", false)
			if err != nil {
				_ = log.Warn("failed to dump current flow counter : ", err)
			}
		}, true)
		if err != nil {
			_ = log.Warn("schedule register failed : ", err)
		}
	}
	//周度流量记录
	if config.Current.Schedule.WeeklyFlowRecord {
		err := service.Register("@weekly", "weekly_flow_recorder", func() {
			storagePath := "./record/weekly/"
			s, err := os.Stat(storagePath)
			if os.IsNotExist(err) || !s.IsDir() {
				_ = os.MkdirAll(storagePath, 0644)
			}
			err = dumpFlowData(storagePath+"weekly_"+time.Now().Format("2006_01_02")+".record", false)
			if err != nil {
				_ = log.Warn("failed to dump current flow counter : ", err)
			}
		}, true)
		if err != nil {
			_ = log.Warn("schedule register failed : ", err)
		}
	}
	//月度流量记录/重置
	if config.Current.Schedule.MonthlyFlowRecord {
		err := service.Register("@monthly", "monthly_flow_recorder", func() {
			storagePath := "./record/monthly/"
			s, err := os.Stat(storagePath)
			if os.IsNotExist(err) || !s.IsDir() {
				_ = os.MkdirAll(storagePath, 0644)
			}
			err = dumpFlowData(storagePath+"monthly_"+time.Now().Format("2006_01_02")+".record", config.Current.Schedule.MonthlyFlowReset)
			if err != nil {
				_ = log.Warn("failed to dump current flow counter : ", err)
			}
		}, true)
		if err != nil {
			_ = log.Warn("schedule register failed : ", err)
		}
	}

}

//
// dumpFlowData
// @Description:
// @param file
//
func dumpFlowData(file string, reset bool) error {
	infos, err := UserServiceInstance().ListUserInfo()
	if err != nil {
		return err
	}
	data := map[string]interface{}{}
	for i := range infos {
		info := model.UserInfo{
			Id: infos[i].Id,
		}
		data[info.Id] = map[string]interface{}{
			"account": infos[i].Account,
			"rx":      infos[i].FlowCount,
			"tx":      infos[i].TXCount,
		}
		if reset {
			err := UserServiceInstance().ResetFlowCounter(&info)
			if err != nil {
				return err
			}
		}
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

//
// checkUserImports
// @Description:
//
func checkUserImports() {
	if UserServiceInstance().configService.UpdatedRecently {
		list, err := UserServiceInstance().infoService.ListFull()
		if err != nil {
			_ = log.Warn("error during user imports check schedule : ", err)
			return
		}
		for i := range list {
			var changed = false
			exports, err := UserServiceInstance().AvailableExports(list[i].Account)
			if err != nil {
				continue
			}
			for j := 0; j < len(list[i].Config.Routes); j++ {
				route := list[i].Config.Routes[j]
				if route.Option != config.RouteOptionImport {
					continue
				}
				if !importValidationCheck(exports, route) {
					changed = true
					log.Info("[", list[i].Account, "][network:", route.Name, "]invalid import : ", route.Network)
					list[i].Config.Routes = append(list[i].Config.Routes[:j], list[i].Config.Routes[j+1:]...)
				}
			}
			if changed {
				_, _ = UserServiceInstance().configService.UpdateById(model.ClientConfig{
					Id:     list[i].ConfigId,
					Routes: list[i].Config.Routes,
					Device: list[i].Config.Device,
					Limit:  list[i].Config.Limit,
				}, false)
			}
		}
		UserServiceInstance().configService.UpdatedRecently = false
	}
}

//
// importValidationCheck
// @Description:
// @param ava
// @param route
// @return bool
//
func importValidationCheck(ava []model.ImportableRoute, route config.Route) bool {
	for i := range ava {
		if ava[i].Network == route.Network {
			return true
		}
	}
	return false
}
