package cache

import (
	"sync"
	"time"
)

//
// IStrStrCache
// @Description:
//
type IStrStrCache interface {
	Get(k string) (v string, ok bool)
	Set(k string, v string)
}

//
// NewIpTable
// @Description:
// @return *IpTable
//
func NewIpTable() *IpTable {
	t := &IpTable{
		expTime: time.Minute * 10,
	}
	t.cleaner()
	return t
}

//
// IpTable
// @Description:
//
type IpTable struct {
	c       sync.Map
	cp      sync.Map
	as      sync.Map
	exp     sync.Map
	expTime time.Duration
}

//
// Get
// @Description:
// @receiver t
// @param k
// @return v
// @return ok
//
func (t *IpTable) Get(k string) (v string, ok bool) {
	if s, ok := t.cp.Load(k); ok {
		return t.ConvertString(s), ok
	}
	if ex, ok := t.exp.Load(k); ok && ex != nil && ex.(int64) <= time.Now().UnixMilli() {
		return "", false
	}
	s, ok := t.c.Load(k)
	//update time
	t.exp.Store(k, time.Now().Add(t.expTime).UnixMilli())
	return t.ConvertString(s), ok
}

//
// Set
// @Description:
// @receiver t
// @param k
// @param v
//
func (t *IpTable) Set(k string, v string) {
	t.c.Store(k, v)
	//update time
	t.exp.Store(k, time.Now().Add(t.expTime).UnixMilli())
}

//
// SetPermanent
// @Description:
// @receiver t
// @param k
// @param v
//
func (t *IpTable) SetPermanent(k string, v string) {
	t.cp.Store(k, v)
	//log.Info("permanent ip table set : ", k, "@", v)
}

//
// GetPermanent
// @Description:
// @receiver t
// @param k
// @return v
// @return ok
//
func (t *IpTable) GetPermanent(k string) (v string, ok bool) {
	s, ok := t.cp.Load(k)
	return t.ConvertString(s), ok
}

//
// AssociatePermanent
// @Description:
// @receiver t
//
func (t *IpTable) AssociatePermanent(a string, k string) {
	if strings, ok := t.as.Load(a); ok && strings != nil {
		t.as.Store(a, append(strings.([]string), k))
	} else {
		t.as.Store(a, []string{k})
	}
}

//
// DelAssociatedPermanent
// @Description:
// @receiver t
// @param a
//
func (t *IpTable) DelAssociatedPermanent(a string) {
	strings, ok := t.as.Load(a)
	t.as.Delete(a)
	if ok && strings != nil {
		for _, s := range strings.([]string) {
			t.DelPermanent(s)
		}
	}
}

//
// DelPermanent
// @Description:
// @receiver t
// @param k
//
func (t *IpTable) DelPermanent(k string) {
	t.cp.Delete(k)
	//log.Info("permanent ip table delete : ", k)
}

//
// cleaner
// @Description:
// @receiver t
//
func (t *IpTable) cleaner() {
	go func() {
		for {
			time.Sleep(time.Minute * 10)
			ts := time.Now().UnixMilli()
			t.exp.Range(func(key, value interface{}) bool {
				if value != nil && value.(int64) <= ts {
					t.exp.Delete(key)
					t.c.Delete(key)
				}
				return true
			})
		}
	}()
}

//
// updateExp
// @Description:
// @receiver t
// @param k
//
func (t *IpTable) updateExp(k string) {
	go func() {
		t.exp.Store(k, time.Now().Add(t.expTime).UnixMilli())
	}()
}

//
// ConvertString
// @Description:
// @receiver t
// @param v
// @return string
//
func (t *IpTable) ConvertString(v interface{}) string {
	if v == nil {
		return ""
	} else {
		return v.(string)
	}
}
