package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

type Redis interface {
	Get(key string) (float64, error)
	Set(key string, value float64, expiration time.Duration) error
}

func NewRedis(addr string) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()

	if _, err := client.Ping(ctx).Result(); err != nil {
		logrus.Fatalf("Failed to connect to Redis: %v", err)
	}

	return &RedisCache{
		client: client,
		ctx:    ctx,
	}
}

func (rc *RedisCache) Get(pair string) (float64, error) {
	val, err := rc.client.Get(rc.ctx, pair).Float64()
	if err == redis.Nil {
		return 0, fmt.Errorf("pair %s not found in cache", pair)
	} else if err != nil {
		return 0, fmt.Errorf("failed to get pair %s from cache: %v", pair, err)
	}
	return val, nil
}

func (rc *RedisCache) Set(pair string, amount float64, expiration time.Duration) error {
	err := rc.client.Set(rc.ctx, pair, amount, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set pair %s in cache: %v", pair, err)
	}
	return nil
}
