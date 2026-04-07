package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"asona/config"
)

// Service defines Redis operations used across the application.
type Service interface {
	// SetSession stores a token → userID mapping with an expiry.
	SetSession(ctx context.Context, token, userID string, ttl time.Duration) error

	// GetUserID returns the userID associated with the given token.
	GetUserID(ctx context.Context, token string) (string, error)

	// CheckLoginSession returns true if the session still exists and belongs to the given userID.
	CheckLoginSession(ctx context.Context, userID, token string) (bool, error)

	// DeleteSession removes the token from the session store.
	DeleteSession(ctx context.Context, token string) error

	// Close closes the Redis connection.
	Close() error

	// Client exposes the underlying *redis.Client for advanced use.
	Client() *redis.Client
}

type service struct {
	client *redis.Client
}

var instance *service

// New returns a singleton Redis Service.
func New() Service {
	if instance != nil {
		return instance
	}

	cfg := config.GetConfig()

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0,
	})

	instance = &service{client: rdb}
	return instance
}

// Client exposes the underlying redis client.
func (s *service) Client() *redis.Client {
	return s.client
}

// Close closes the Redis client connection.
func (s *service) Close() error {
	return s.client.Close()
}
