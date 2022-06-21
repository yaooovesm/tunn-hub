package administration

import (
	"time"
	"tunn-hub/config"
	"tunn-hub/monitor"
)

//
// recorderService
// @Description:
//
type recorderService struct {
	TrafficRecorder monitor.HubTrafficRecorder
}

//
// newRecorderService
// @Description:
// @return *recorderService
//
func newRecorderService() *recorderService {
	return &recorderService{
		TrafficRecorder: monitor.NewHubTrafficRecorder(),
	}
}

//
// RecentHour
// @Description:
// @receiver rs
// @return []monitor.HubTrafficStamp
//
func (rs *recorderService) RecentHour() []monitor.HubTrafficStamp {
	return rs.TrafficRecorder.Recent(60)
}

//
// RecentDay
// @Description:
// @receiver rs
// @return []monitor.HubTrafficStamp
//
func (rs *recorderService) RecentDay() []monitor.HubTrafficStamp {
	data := make([]monitor.HubTrafficStamp, 0)
	path := "./record/hourly/"
	if config.Current.Schedule.HourlyNetworkTrafficDump != "" {
		path = config.Current.Schedule.HourlyNetworkTrafficDump
	}
	now := time.Now().UnixMilli()
	duration := (time.Hour * time.Duration(24)).Milliseconds()
	for i := 0; i <= 23; i++ {
		history := readTrafficHistory(path, i, now, duration)
		data = append(data, history...)
	}
	data = append(data, rs.TrafficRecorder.Latest()...)
	sortShell(data)
	return data
}

//
// sortShell
// @Description:
// @param list
//
func sortShell(list []monitor.HubTrafficStamp) {
	// 数组长度
	n := len(list)
	// 每次减半，直到步长为 1
	for step := n / 2; step >= 1; step /= 2 {
		// 开始插入排序，每一轮的步长为 step
		for i := step; i < n; i += step {
			for j := i - step; j >= 0; j -= step {
				// 满足插入那么交换元素
				if list[j+step].Timestamp < list[j].Timestamp {
					list[j], list[j+step] = list[j+step], list[j]
					continue
				}
				break
			}
		}
	}
}
