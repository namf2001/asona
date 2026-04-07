package redis

import (
	"context"
	"fmt"
	"time"
)

const sessionKeyPrefix = "session:"

// sessionKey builds the Redis key for a given token.
func sessionKey(token string) string {
	return fmt.Sprintf("%s%s", sessionKeyPrefix, token)
}

// SetSession stores token → userID with a TTL.
func (s *service) SetSession(ctx context.Context, token, userID string, ttl time.Duration) error {
	return s.client.Set(ctx, sessionKey(token), userID, ttl).Err()
}

// GetUserID retrieves the userID for the given token.
// Returns ErrSessionNotFound if the token does not exist.
func (s *service) GetUserID(ctx context.Context, token string) (string, error) {
	userID, err := s.client.Get(ctx, sessionKey(token)).Result()
	if err != nil {
		return "", ErrSessionNotFound
	}
	return userID, nil
}

// CheckLoginSession verifies that the token exists and belongs to the given userID.
// Returns false (not an error) if token is expired or belongs to a different user.
func (s *service) CheckLoginSession(ctx context.Context, userID, token string) (bool, error) {
	storedUserID, err := s.GetUserID(ctx, token)
	if err != nil {
		return false, nil // session not found → expired/invalid
	}
	return storedUserID == userID, nil
}

// DeleteSession removes a token from Redis (logout).
func (s *service) DeleteSession(ctx context.Context, token string) error {
	return s.client.Del(ctx, sessionKey(token)).Err()
}
