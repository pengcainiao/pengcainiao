package redis

import (
	"context"
	"time"
)

type RedisCache struct {
	RedisNode
}

func (r RedisCache) Get(key string) interface{} {
	value, err := r.RedisNode.Get(context.Background(), key).Result()
	if err != nil {
		return nil
	}
	return value
}

func (r RedisCache) Set(key string, val interface{}, timeout time.Duration) error {
	_, err := r.RedisNode.Set(context.Background(), key, val, timeout).Result()
	return err
}

func (r RedisCache) IsExist(key string) bool {
	ic, err := r.RedisNode.Exists(context.Background(), key).Result()
	if err != nil {
		return false
	}
	return ic > 0
}

func (r RedisCache) Delete(key string) error {
	_, err := r.RedisNode.Del(context.Background(), key).Result()
	return err
}
