package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// Publish publishes a message to the specified channel using WebSocket client.
func (s service) Publish(ctx context.Context, channel string, message interface{}) error {
	return s.wsClient.Publish(ctx, channel, message).Err()
}

// Subscribe subscribes to the specified channels using WebSocket client.
func (s service) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return s.wsClient.Subscribe(ctx, channels...)
}
