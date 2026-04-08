package redis

import (
	"context"
	"time"
)

// Set stores a key-value pair with optional TTL.
func (s service) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return s.client.Set(ctx, key, value, ttl).Err()
}

// Get retrieves the value for the given key.
func (s service) Get(ctx context.Context, key string) (string, error) {
	return s.client.Get(ctx, key).Result()
}

// Del removes the specified keys.
func (s service) Del(ctx context.Context, keys ...string) error {
	return s.client.Del(ctx, keys...).Err()
}

// Exists checks if the specified keys exist.
func (s service) Exists(ctx context.Context, keys ...string) (int64, error) {
	return s.client.Exists(ctx, keys...).Result()
}
