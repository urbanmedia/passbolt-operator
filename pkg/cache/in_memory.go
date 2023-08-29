package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var (
	// Ensure inMemory implements Cacher.
	_ Cacher = (*inMemory)(nil)
)

type inMemory struct {
	// mu is used to prevent concurrent access to the secret cache.
	mu        sync.RWMutex
	cacheData map[string]any
}

// NewInMemoryCache returns a new in-memory cache.
func NewInMemoryCache() Cacher {
	return &inMemory{
		cacheData: make(map[string]any),
	}
}

// Set sets the value for a key with a TTL.
// ttl is ignored for in-memory cache.
func (c *inMemory) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cacheData[key] = value
	return nil
}

// Get returns the value for a key.
func (c *inMemory) Get(ctx context.Context, key string) (any, error) {
	value, ok := c.cacheData[key]
	if !ok {
		return nil, fmt.Errorf("cache miss for key %q", key)
	}
	return value, nil
}

// Delete deletes the value for a key.
// If the key does not exist, Delete returns nil (no error).
func (c *inMemory) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cacheData, key)
	return nil
}
