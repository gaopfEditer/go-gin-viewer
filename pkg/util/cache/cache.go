package cache

import (
	"time"
)

const Expire = 3600 * time.Second

type Cache interface {
	Set(key string, value interface{}, expire ...time.Duration) error
	Get(key string) (interface{}, error)
	Del(key string) error
	DelByPrefix(key string) error
	Update(key, value interface{}) error
}
