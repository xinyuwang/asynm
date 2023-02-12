package asynm

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisClient struct {
	client redis.UniversalClient
}

func newRedisClient(RedisCacheConfig *redis.UniversalOptions) (*redisClient, error) {

	client := redis.NewUniversalClient(RedisCacheConfig)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis.ping error: %w", err)
	}

	return &redisClient{
		client: client,
	}, nil
}

func (r *redisClient) Expire(key string, expiration time.Duration) error {
	return r.client.Expire(context.Background(), key, expiration).Err()
}

func (r *redisClient) ExpireAt(key string, tm time.Time) error {
	return r.client.ExpireAt(context.Background(), key, tm).Err()
}

func (r *redisClient) HSet(key string, pairs ...interface{}) error {
	return r.client.HSet(context.Background(), key, pairs...).Err()
}

func (r *redisClient) HGet(key string, field string) string {
	return r.client.HGet(context.Background(), key, field).Val()
}

func (r *redisClient) HGetAll(key string) map[string]string {
	return r.client.HGetAll(context.Background(), key).Val()
}

func (r *redisClient) HIncrBy(key string, field string, incr int64) error {
	return r.client.HIncrBy(context.Background(), key, field, incr).Err()
}
