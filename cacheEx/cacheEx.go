package cacheEx

import (
	"time"
)

var (
	cacheMap = make(map[string]*cacheItem)
	defExp   = time.Hour * 1
)

type cacheItem struct {
	expTime time.Time
	expDur  time.Duration
}

func init() {
	go run()
}

func Set(key string, exp time.Duration) {
	if exp == 0 {
		exp = defExp
	}
	cacheMap[key] = &cacheItem{
		expTime: time.Now().Add(exp),
		expDur:  exp,
	}
}

func Get(key string) (*cacheItem, bool) {
	it, b := cacheMap[key]
	return it, b
}

func Refresh(key string) {
	if it, b := cacheMap[key]; b {
		Set(key, it.expDur)
	}
}

func Del(key string) {
	delete(cacheMap, key)
}

func run() {
	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-ticker.C:
			for k, it := range cacheMap {
				if it.expTime.Before(time.Now()) {
					Del(k)
				}
			}
		}
	}
}
