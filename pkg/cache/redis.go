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
func NewRedisCache(logr logr.Logger, config *redis.Options) Cacher {
	rdb := redis.NewClient(config)
	return &Redis{
		logr:   logr.WithName("redis"),
		client: rdb,
	}
}

// Close closes the redis connection.
func (r *Redis) Close() error {
	return r.client.Close()
}

func (r *Redis) Get(ctx context.Context, key string) (any, error) {
	err := r.client.Get(ctx, key).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to read key %q: %w", key, err)
	}
	return nil, nil
}

func (r *Redis) Set(ctx context.Context, key string, val any, ttl time.Duration) error {
	err := r.client.Set(ctx, key, val, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to write key %q: %w", key, err)
	}
	return nil
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key %q: %w", key, err)
	}
	return nil
}
