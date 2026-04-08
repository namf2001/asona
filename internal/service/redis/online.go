package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"asona/internal/constants"
)

// SetUserOnline marks a user as online with a TTL.
func (s service) SetUserOnline(ctx context.Context, userID string, ttl time.Duration) error {
	key := constants.UserOnlineKeyPrefix + userID
	return s.client.Set(ctx, key, "1", ttl).Err()
}

// SetUserOffline removes the user's online status.
func (s service) SetUserOffline(ctx context.Context, userID string) error {
	key := constants.UserOnlineKeyPrefix + userID
	return s.client.Del(ctx, key).Err()
}

// IsUserOnline checks if a user is currently online.
func (s service) IsUserOnline(ctx context.Context, userID string) (bool, error) {
	key := constants.UserOnlineKeyPrefix + userID
	result, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// GetOnlineUsers returns a list of online users from the given user IDs.
func (s service) GetOnlineUsers(ctx context.Context, userIDs []string) ([]string, error) {
	if len(userIDs) == 0 {
		return []string{}, nil
	}

	keys := make([]string, len(userIDs))
	for i, id := range userIDs {
		keys[i] = constants.UserOnlineKeyPrefix + id
	}

	pipe := s.client.Pipeline()
	cmds := make([]*redis.IntCmd, len(keys))
	for i, key := range keys {
		cmds[i] = pipe.Exists(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	onlineUsers := make([]string, 0)
	for i, cmd := range cmds {
		if cmd.Val() > 0 {
			onlineUsers = append(onlineUsers, userIDs[i])
		}
	}

	return onlineUsers, nil
}
