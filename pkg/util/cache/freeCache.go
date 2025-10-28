package cache

import (
	"cambridge-hit.com/gin-base/activateserver/pkg/util/logger"
	"cambridge-hit.com/gin-base/activateserver/resource"
	"errors"
	"github.com/coocood/freecache"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"strings"
	"time"
)

var MyFreeCache *FreeCache

type FreeCache struct {
	fch *freecache.Cache
}

func NewFreeCache() *FreeCache {
	cacheSize := 100 * 1024 * 1024 // 缓存大小，单位为字节
	cache := freecache.NewCache(cacheSize)
	return &FreeCache{fch: cache}
}

func (f *FreeCache) Set(key string, val interface{}, expire ...time.Duration) error {
	eTime := Expire
	if expire != nil && len(expire) > 0 {
		eTime = expire[0]
	}
	var cache = make([]byte, 0)
	var err error
	// 如果value是[]byte就不需要转化了
	switch val.(type) {
	case []byte:
		cache = val.([]byte)
	default:
		cache, err = jsonEncode(val)
	}
	if err == nil {
		err = f.fch.Set([]byte(key), cache, int(eTime))
	}
	return err
}

func (f *FreeCache) Get(key string) (interface{}, error) {
	cache, err := f.fch.Get([]byte(key))
	if err != nil {
		return nil, err
	}
	var result interface{}
	err = jsonDecode(cache, &result)
	return result, err
}

func (f *FreeCache) Del(key string) error {
	affected := f.fch.Del([]byte(key))
	if !affected {
		return errors.New("key not found")
	}
	return nil
}

func (f *FreeCache) Update(key, value string) {
	// Check if the key exists
	_, err := f.fch.Get([]byte(key))
	if err != nil {
		logger.Error("freecache error:", zap.Error(err))
		return
	}
	err = f.Set(key, value, time.Duration(resource.Conf.AccessExpire)*time.Second)
	if err != nil {
		logger.Error("freecache error:", zap.Error(err))
		return
	}
}

func (f *FreeCache) DelByPrefix(pre string) error {
	iter := f.fch.NewIterator()
	for {
		entry := iter.Next()
		if entry == nil {
			break
		}
		key := string(entry.Key)
		if strings.HasPrefix(key, pre) {
			f.fch.Del(entry.Key)
		}
	}
	return nil
}

// json序列化
func jsonEncode(data interface{}) ([]byte, error) {
	enc, err := jsoniter.Marshal(data)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

// json反序列化
func jsonDecode(data []byte, to interface{}) error {
	err := jsoniter.Unmarshal(data, to)
	return err
}
