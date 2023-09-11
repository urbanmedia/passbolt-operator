package cache

import (
	"context"
	"time"
)

// Cacher is the interface that wraps the basic Set and Get methods for a cache.
type Cacher interface {
	// Set sets the value for a key with a TTL.
	// If the TTL is 0, the key will never expire.
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	// Get returns the value for a key.
	Get(ctx context.Context, key string) ([]byte, error)
	// Delete deletes the value for a key.
	Delete(ctx context.Context, key string) error
}
