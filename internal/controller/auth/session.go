package auth

import (
	"context"
	"fmt"
	"time"

	"asona/config"
	"asona/internal/model"
	"asona/internal/pkg/jwt"

	pkgerrors "github.com/pkg/errors"
)

// issueSession creates a JWT token and persists the backing session row.
func (i impl) issueSession(ctx context.Context, user model.User, userAgent, ipAddress string) (string, error) {
	token, err := jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		return "", pkgerrors.WithStack(fmt.Errorf("failed to generate token: %w", err))
	}

	cfg := config.GetConfig()
	accessDuration := cfg.JWTAccessDuration
	if accessDuration == 0 {
		accessDuration = time.Hour
	}

	_, err = i.repo.Session().Create(ctx, model.Session{
		UserID:       user.ID,
		SessionToken: token,
		ExpiresAt:    time.Now().Add(accessDuration),
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
	})
	if err != nil {
		return "", pkgerrors.WithStack(fmt.Errorf("failed to persist session: %w", err))
	}

	return token, nil
}
