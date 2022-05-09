package cache

import (
	"sync"
)

//
// IpTableV2Permanent
// @Description:
//
type IpTableV2Permanent struct {
	lock sync.Mutex
	m    map[string]string
}

//
// NewIpTableV2P
// @Description:
// @param clearInterval
// @param expire
//
func NewIpTableV2P() *IpTableV2Permanent {
	return &IpTableV2Permanent{
		lock: sync.Mutex{},
		m:    make(map[string]string),
	}
}

//
// Set
// @Description:
// @receiver ipt
// @param k
// @param v
//
func (ipt *IpTableV2Permanent) Set(k string, v string) {
	ipt.lock.Lock()
	ipt.m[k] = v
	ipt.lock.Unlock()
}

//
// Get
// @Description:
// @receiver ipt
// @param k
// @return v
// @return ok
//
func (ipt *IpTableV2Permanent) Get(k string) (v string, ok bool) {
	ipt.lock.Lock()
	v, ok = ipt.m[k]
	ipt.lock.Unlock()
	return
}
