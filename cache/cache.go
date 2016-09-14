package cache

import (
	"time"
)

var instance ICache

func RegisterCacheProvider(s ICache) {
	instance = s
}

func GetCacheInstance() ICache {
	return instance
}

func Set(key string, obj interface{}, d time.Duration) error {
	if instance == nil {
		panic("cache instance is nil")
	}
	return instance.Set(key, obj, d)
}

func Get(key string, obj interface{}) error {
	if instance == nil {
		panic("cache instance is nil")
	}
	return instance.Get(key, obj)
}

func Delete(key string) error {
	if instance == nil {
		panic("cache instance is nil")
	}
	return instance.Delete(key)
}
