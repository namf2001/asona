package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"asona/config"
	"asona/internal/constants"
)

// Service defines Redis operations used across the application.
type Service interface {
	// Session operations
	SetSession(ctx context.Context, token, userID string, ttl time.Duration) error
	GetUserID(ctx context.Context, token string) (string, error)
	CheckLoginSession(ctx context.Context, userID, token string) (bool, error)
	DeleteSession(ctx context.Context, token string) error

	// Pub/Sub operations for WebSocket
	Publish(ctx context.Context, channel string, message interface{}) error
	Subscribe(ctx context.Context, channels ...string) *redis.PubSub

	// User online status operations
	SetUserOnline(ctx context.Context, userID string, ttl time.Duration) error
	SetUserOffline(ctx context.Context, userID string) error
	IsUserOnline(ctx context.Context, userID string) (bool, error)
	GetOnlineUsers(ctx context.Context, userIDs []string) ([]string, error)

	// General cache operations
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, keys ...string) (int64, error)

	// Connection management
	Close() error
	Client() *redis.Client
	WSClient() *redis.Client
}

type service struct {
	client   *redis.Client // For session/general operations
	wsClient *redis.Client // For WebSocket pub/sub
}

// New returns a Redis Service. No longer a pointer-based singleton to avoid nil errors.
func New() Service {
	cfg := config.GetConfig()
	addr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)

	// Main client for session operations
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
		DB:       constants.RedisDBSession,
	})

	// WebSocket client for pub/sub
	wsClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
		DB:       constants.RedisDBWebSocket,
	})

	return service{
		client:   client,
		wsClient: wsClient,
	}
}

// Client exposes the underlying redis client for session operations.
func (s service) Client() *redis.Client {
	return s.client
}

// WSClient exposes the WebSocket redis client.
func (s service) WSClient() *redis.Client {
	return s.wsClient
}

// Close closes all Redis client connections.
func (s service) Close() error {
	if err := s.client.Close(); err != nil {
		return err
	}
	return s.wsClient.Close()
}
