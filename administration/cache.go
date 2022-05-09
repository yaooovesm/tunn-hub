package administration

import (
	"sync"
	"time"
)

const (
	DefaultTimeout = 10 * 60 * 1000
	CleanInterval  = 1 * 60 * 1000
)

//
// Cache
// @Description:
//
type Cache struct {
	m              sync.Map
	exp            sync.Map
	defaultTimeout int64
}

//
// NewCache
// @Description:
// @param defaultTimeout
// @return *Cache
//
func NewCache(defaultTimeout int64) *Cache {
	c := &Cache{
		m:              sync.Map{},
		exp:            sync.Map{},
		defaultTimeout: defaultTimeout,
	}
	if c.defaultTimeout != 0 {
		c.autoClean()
	}
	return c
}

//
// SetWithTimeout
// @Description:
// @receiver c
// @param key
// @param value
// @param timeout
//
func (c *Cache) SetWithTimeout(key string, value interface{}, timeout int64) {
	c.m.Store(key, value)
	c.exp.Store(key, time.Now().UnixMilli()+timeout)
}

//
// Set
// @Description:
// @receiver c
// @param key
// @param value
//
func (c *Cache) Set(key string, value interface{}) {
	c.m.Store(key, value)
	if c.defaultTimeout != 0 {
		c.exp.Store(key, time.Now().UnixMilli()+c.defaultTimeout)
	}
}

//
// Get
// @Description:
// @receiver c
// @param key
// @return interface{}
//
func (c *Cache) Get(key string) interface{} {
	if value, ok := c.m.Load(key); ok {
		return value
	}
	return nil
}

//
// Delete
// @Description:
// @receiver c
//
func (c *Cache) Delete(key string) {
	c.m.Delete(key)
	c.exp.Delete(key)
}

//
// autoClean
// @Description:
// @receiver c
//
func (c *Cache) autoClean() {
	go func() {
		interval := time.Millisecond * time.Duration(CleanInterval)
		for {
			current := time.Now().UnixMilli()
			c.exp.Range(func(key, value interface{}) bool {
				if value.(int64) < current {
					c.m.Delete(key)
					c.exp.Delete(key)
				}
				return true
			})
			time.Sleep(interval)
		}
	}()
}
