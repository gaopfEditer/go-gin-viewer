package cache

import (
	"cambridge-hit.com/gin-base/activateserver/resource"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var MyRedis *RedisCache

// RedisCache 实现了 Cache 接口
type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisCache 创建一个 RedisCache 实例
func NewRedisCache() *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     resource.Conf.RedisConfig.Addr,
		Password: resource.Conf.RedisConfig.Password,
		DB:       resource.Conf.RedisConfig.DB,
	})
	return &RedisCache{
		client: rdb,
		ctx:    context.Background(),
	}
}

// Set 实现 Cache 接口中的 Set 方法
func (c *RedisCache) Set(key string, value interface{}, expire ...time.Duration) error {
	var exp time.Duration
	if len(expire) > 0 {
		exp = expire[0]
	}
	return c.client.Set(c.ctx, key, value, exp).Err()
}

// Get 实现 Cache 接口中的 Get 方法
func (c *RedisCache) Get(key string) (interface{}, error) {
	return c.client.Get(c.ctx, key).Result()
}

// Del 实现 Cache 接口中的 Del 方法
func (c *RedisCache) Del(key string) error {
	return c.client.Del(c.ctx, key).Err()
}

// DelByPrefix 实现 Cache 接口中的 DelByPrefix 方法
func (c *RedisCache) DelByPrefix(prefix string) error {
	keys, err := c.client.Keys(c.ctx, prefix+"*").Result()
	if err != nil {
		return err
	}
	for _, key := range keys {
		if err := c.client.Del(c.ctx, key).Err(); err != nil {
			return err
		}
	}
	return nil
}

// Update 方法
func (c *RedisCache) Update(key string, value interface{}) error {
	// 获取键的剩余过期时间
	ttl, err := c.client.TTL(c.ctx, key).Result()
	if err != nil {
		return err
	}
	if ttl <= 0 {
		// 键不存在或没有过期时间
		return c.client.Set(c.ctx, key, value, 0).Err()
	}

	// 使用剩余的过期时间设置新的值
	return c.client.Set(c.ctx, key, value, ttl).Err()
}

// HSet 实现 Cache 接口中的 HSet 方法
func (c *RedisCache) HSet(key string, value interface{}, expire ...time.Duration) error {
	ctx := context.Background()

	// 将Map存储到Redis中
	err := c.client.HSet(ctx, key, value).Err()
	if err != nil {
		return err
	}

	// 如果有过期时间参数，设置过期时间
	if len(expire) > 0 {
		err = c.client.Expire(ctx, key, expire[0]).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

// HGet 实现 Cache 接口中的 HGet 方法
func (c *RedisCache) HGet(key string, field string) (string, error) {
	ctx := context.Background()

	// 从Redis中获取指定字段的值
	value, err := c.client.HGet(ctx, key, field).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", fmt.Errorf("field %s does not exist in key %s", field, key)
		}
		return "", err
	}

	return value, nil
}

// HDel 实现 Cache 接口中的 HDel 方法
func (c *RedisCache) HDel(key string, fields ...string) error {
	ctx := context.Background()

	// 从Redis中删除指定字段
	err := c.client.HDel(ctx, key, fields...).Err()
	if err != nil {
		return err
	}

	return nil
}
