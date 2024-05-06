package cache

import (
	"sync"
	"time"
)

var (
	tickerMap = make(map[string]*time.Ticker)
	defExp    = time.Hour * 1
)

type Cache struct {
	mu   sync.RWMutex
	data map[string]cacheItem
}

type cacheItem struct {
	items   interface{}
	expTime time.Time
	expDur  time.Duration
}

func New() *Cache {
	c := &Cache{
		data: make(map[string]cacheItem),
	}
	return c
}

func (c *Cache) Set(key string, value interface{}, expiration ...time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, b := tickerMap[key]; b {
		tickerMap[key].Stop()
		delete(tickerMap, key)
	}
	expireTime := time.Now().Add(defExp)
	var expDur time.Duration
	if expiration != nil {
		expDur = expiration[0]
		expireTime = time.Now().Add(expDur)
	}
	c.data[key] = cacheItem{
		items:   value,
		expTime: expireTime,
		expDur:  expDur,
	}
	go c.run(key)
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.data[key]
	if ok && item.expTime.After(time.Now()) {
		return item.items, true
	}
	return nil, false
}

func (c *Cache) Del(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

func (c *Cache) run(key string) {
	tickerMap[key] = time.NewTicker(2 * time.Second)
	for {
		if _, b := tickerMap[key]; !b {
			return
		}
		select {
		case <-tickerMap[key].C:
			it, exist := c.data[key]
			if exist {
				if it.expTime.Before(time.Now()) {
					c.Del(key)
					tickerMap[key].Stop()
					delete(tickerMap, key)
					return
				}
			} else {
				tickerMap[key].Stop()
				delete(tickerMap, key)
				return
			}
		}
	}
}
