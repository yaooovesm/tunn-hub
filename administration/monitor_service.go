package administration

import (
	"time"
	"tunn-hub/monitor"
)

//
// monitorService
// @Description:
//
type monitorService struct {
	data *monitor.SystemData
	gap  int
}

//
// newMonitorService
// @Description:
// @return *monitorService
//
func newMonitorService(gap int) *monitorService {
	serv := &monitorService{data: &monitor.SystemData{}, gap: gap}
	serv.data.CollectAll()
	go serv.autoCollect()
	return serv
}

//
// autoCollect
// @Description:
// @receiver m
//
func (m *monitorService) autoCollect() {
	for {
		time.Sleep(time.Millisecond * time.Duration(m.gap))
		m.data.CollectAll()
	}
}

//
// GetSystemData
// @Description:
// @receiver m
//
func (m *monitorService) GetSystemData() monitor.SystemData {
	return *m.data
}
