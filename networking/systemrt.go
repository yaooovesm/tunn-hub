package networking

import (
	"sync"
	"tunn-hub/config"
)

//
// RouteRecord
// @Description:
//
type RouteRecord struct {
	DeviceName string
	Network    string
	Deployed   bool
	Err        error
}

//
// Deploy
// @Description:
// @receiver r
//
func (r *RouteRecord) Deploy() {
	if r.Deployed {
		return
	}
	err := AddSystemRoute(r.Network, r.DeviceName)
	if err != nil {
		r.Deployed = false
		r.Err = err
		return
	} else {
		r.Deployed = true
	}
}

//
// SystemRouteTable
// @Description:
//
type SystemRouteTable struct {
	DeviceName string
	Records    map[string]*RouteRecord
	lock       sync.RWMutex
}

//
// NewSystemRouteTable
// @Description:
// @return *SystemRouteTable
//
func NewSystemRouteTable(devName string) *SystemRouteTable {
	return &SystemRouteTable{
		DeviceName: devName,
		Records:    make(map[string]*RouteRecord),
		lock:       sync.RWMutex{},
	}
}

//
// Merge
// @Description:
// @receiver t
// @param routes
//
func (t *SystemRouteTable) Merge(routes []config.Route) {
	t.lock.Lock()
	defer t.lock.Unlock()
	for i := range routes {
		route := routes[i]
		//只导入import
		if route.Option == config.RouteOptionImport {
			//检查是否已导入
			if _, ok := t.Records[route.Network]; !ok {
				//导入记录
				t.Records[route.Network] = &RouteRecord{
					DeviceName: t.DeviceName,
					Network:    route.Network,
					Deployed:   false,
					Err:        nil,
				}
			}
		}
	}
}

//
// DeployAll
// @Description:
// @receiver t
//
func (t *SystemRouteTable) DeployAll() {
	t.lock.Lock()
	defer t.lock.Unlock()
	for n := range t.Records {
		t.Records[n].Deploy()
	}
}
