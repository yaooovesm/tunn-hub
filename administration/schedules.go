package administration

import (
	"encoding/json"
	log "github.com/cihub/seelog"
	"io/ioutil"
	"os"
	"time"
	"tunn-hub/administration/model"
	"tunn-hub/config"
)

//
// setupSchedules
// @Description:
// @param service
//
func setupSchedules(service *scheduleService) {
	//err := service.Register("@every 30s", "flow_recorder", func() {
	//	if ServerServiceInstance() == nil {
	//		return
	//	}
	//	rx := uint64(0)
	//	tx := uint64(0)
	//	if ServerServiceInstance().rxFlowCounter != nil {
	//		rx = ServerServiceInstance().rxFlowCounter.FlowSpeed
	//	}
	//	if ServerServiceInstance().txFlowCounter != nil {
	//		tx = ServerServiceInstance().txFlowCounter.FlowSpeed
	//	}
	//	fmt.Println()
	//	fmt.Println("-------------------------")
	//	fmt.Println("rx <-- ", rx)
	//	fmt.Println("tx <-- ", tx)
	//	fmt.Println("-------------------------")
	//	fmt.Println()
	//}, true)
	//if err != nil {
	//	_ = log.Warn("flow_recorder register failed : ", err)
	//	return
	//}
	//daily recorder
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
	//weekly recorder
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
	//monthly recorder
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
