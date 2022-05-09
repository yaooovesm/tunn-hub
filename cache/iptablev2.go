package cache

import (
	"sync"
	"time"
)

//
// IpTableV2
// @Description:
//
type IpTableV2 struct {
	cleanInterval time.Duration
	expire        time.Duration
	lock          sync.Mutex
	m             map[string]string
	exp           map[string]int64
}

//
// NewIpTableV2
// @Description:
// @param clearInterval
// @param expire
//
func NewIpTableV2(clearInterval time.Duration, expire time.Duration) *IpTableV2 {
	return &IpTableV2{
		cleanInterval: clearInterval,
		expire:        expire,
		lock:          sync.Mutex{},
		m:             make(map[string]string),
		exp:           make(map[string]int64),
	}
}

//
// Set
// @Description:
// @receiver ipt
// @param k
// @param v
//
func (ipt *IpTableV2) Set(k string, v string) {
	ipt.lock.Lock()
	ipt.m[k] = v
	ipt.exp[k] = time.Now().Add(ipt.expire).UnixMilli()
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
func (ipt *IpTableV2) Get(k string) (v string, ok bool) {
	ipt.lock.Lock()
	v, ok = ipt.m[k]
	ipt.lock.Unlock()
	return
}

//
// autoCleaner
// @Description:
// @receiver ipt
//
func (ipt *IpTableV2) autoCleaner() {
	go func() {
		for {
			time.Sleep(ipt.cleanInterval)
			ts := time.Now().UnixMilli()
			for k := range ipt.exp {
				if exp, ok := ipt.exp[k]; ok && exp >= ts {
					ipt.lock.Lock()
					delete(ipt.exp, k)
					delete(ipt.m, k)
					ipt.lock.Unlock()
				}
			}
		}
	}()
}
