package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/redis/go-redis/v9"
)

var (
	// Ensure redis implements Cacher.
	_ Cacher = (*Redis)(nil)
)

type Redis struct {
	logr   logr.Logger
	client *redis.Client
}

// NewRedisCache returns a new redis cache.
func NewRedisCache(logr logr.Logger, config *redis.Options) *Redis {
	myLog := logr.WithName("redis")
	myLog.Info("establishing redis connection...")
	rdb := redis.NewClient(config)
	return &Redis{
		logr:   myLog,
		client: rdb,
	}
}

func (r *Redis) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Close closes the redis connection.
func (r *Redis) Close() error {
	r.logr.Info("closing redis connection...")
	return r.client.Close()
}

func (r *Redis) Get(ctx context.Context, key string) (any, error) {
	r.logr.Info("reading key from redis", "key", key)
	rsp := r.client.Get(ctx, key)
	if err := rsp.Err(); err != nil {
		return nil, fmt.Errorf("failed to read key %q: %w", key, err)
	}
	return rsp.Val(), nil
}

func (r *Redis) Set(ctx context.Context, key string, val any, ttl time.Duration) error {
	r.logr.Info("writing key to redis", "key", key)
	err := r.client.Set(ctx, key, val, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to write key %q: %w", key, err)
	}
	return nil
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	r.logr.Info("deleting key from redis", "key", key)
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key %q: %w", key, err)
	}
	return nil
}
