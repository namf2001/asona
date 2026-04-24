package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"asona/internal/model"
	"asona/internal/repository"
	"asona/internal/repository/accounts"
	"asona/internal/repository/users"

	pkgerrors "github.com/pkg/errors"
)

const googleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"

type googleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

// GoogleAuthURL returns the Google authorization URL for the given OAuth state.
func (i impl) GoogleAuthURL(ctx context.Context, state string) (string, error) {
	cfg := i.oauth.GoogleConfig()
	if cfg == nil {
		return "", pkgerrors.WithStack(ErrOAuthNotConfigured)
	}

	_ = ctx
	return cfg.AuthCodeURL(state), nil
}

// GoogleCallback exchanges the OAuth code, loads the Google profile, and issues a local session.
func (i impl) GoogleCallback(ctx context.Context, code string) (model.User, string, error) {
	cfg := i.oauth.GoogleConfig()
	if cfg == nil {
		return model.User{}, "", pkgerrors.WithStack(ErrOAuthNotConfigured)
	}

	token, err := cfg.Exchange(ctx, code)
	if err != nil {
		return model.User{}, "", pkgerrors.WithStack(fmt.Errorf("%w: %v", ErrOAuthExchangeFailed, err))
	}

	profile, err := fetchGoogleUserInfo(ctx, token.AccessToken)
	if err != nil {
		return model.User{}, "", err
	}
	if profile.Email == "" {
		return model.User{}, "", pkgerrors.WithStack(fmt.Errorf("google oauth profile missing email"))
	}
	if !profile.VerifiedEmail {
		return model.User{}, "", pkgerrors.WithStack(ErrOAuthEmailNotVerified)
	}

	var user model.User
	var sessionToken string

	err = i.repo.DoInTx(ctx, func(ctx context.Context, txRepo repository.Registry) error {
		account, err := txRepo.Account().GetByProvider(ctx, "google", profile.ID)
		if err == nil {
			// Existing Google account — load user and issue session.
			user, err = txRepo.User().GetByID(ctx, account.UserID)
			if err != nil {
				return err
			}
		} else {
			if !errors.Is(err, accounts.ErrAccountNotFound) {
				return err
			}

			// No Google account yet — find or create the user row.
			user, err = txRepo.User().GetByEmail(ctx, profile.Email)
			if err != nil {
				if !errors.Is(err, users.ErrUserNotFound) {
					return err
				}
				user, err = createGoogleUser(ctx, txRepo, profile)
				if err != nil {
					return err
				}
			}

			_, err = txRepo.Account().Create(ctx, model.Account{
				UserID:            user.ID,
				Provider:          "google",
				ProviderAccountID: profile.ID,
				AccessToken:       token.AccessToken,
				RefreshToken:      token.RefreshToken,
				IDToken:           tokenValue(token.Extra("id_token")),
				Scope:             tokenValue(token.Extra("scope")),
			})
			if err != nil {
				if errors.Is(err, accounts.ErrAccountAlreadyExists) {
					existingAccount, lookupErr := txRepo.Account().GetByProvider(ctx, "google", profile.ID)
					if lookupErr != nil {
						return lookupErr
					}
					user, lookupErr = txRepo.User().GetByID(ctx, existingAccount.UserID)
					if lookupErr != nil {
						return lookupErr
					}
				} else {
					return err
				}
			}
		}

		return nil
	}, nil)
	if err != nil {
		return model.User{}, "", err
	}

	// Issue session outside the transaction using the standard helper
	sessionToken, err = i.issueSession(ctx, user, "", "")
	if err != nil {
		return model.User{}, "", err
	}

	return user, sessionToken, nil
}

func fetchGoogleUserInfo(ctx context.Context, accessToken string) (googleUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, googleUserInfoURL, nil)
	if err != nil {
		return googleUserInfo{}, pkgerrors.WithStack(fmt.Errorf("failed to create google userinfo request: %w", err))
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return googleUserInfo{}, pkgerrors.WithStack(fmt.Errorf("%w: %v", ErrOAuthUserInfoFailed, err))
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return googleUserInfo{}, pkgerrors.WithStack(fmt.Errorf("%w: unexpected status %d", ErrOAuthUserInfoFailed, resp.StatusCode))
	}

	var profile googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return googleUserInfo{}, pkgerrors.WithStack(fmt.Errorf("failed to decode google userinfo: %w", err))
	}

	return profile, nil
}

func createGoogleUser(ctx context.Context, txRepo repository.Registry, profile googleUserInfo) (model.User, error) {
	parts := strings.Split(profile.Email, "@")
	prefix := strings.ToLower(parts[0])
	if prefix == "" {
		prefix = "google-user"
	}
	username := fmt.Sprintf("%s-%s", sanitizeUsername(prefix), shortAccountSuffix(profile.ID))

	user := model.User{
		Name:     profile.Name,
		Username: username,
		Email:    profile.Email,
		Image:    profile.Picture,
	}

	created, err := txRepo.User().Create(ctx, user)
	if err != nil {
		return model.User{}, err
	}

	return created, nil
}

func sanitizeUsername(value string) string {
	var b strings.Builder
	b.Grow(len(value))
	for _, r := range strings.ToLower(value) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			b.WriteRune(r)
			continue
		}
		if r == '.' || r == ' ' || r == '+' {
			b.WriteRune('-')
		}
	}
	result := strings.Trim(b.String(), "-_")
	if result == "" {
		return "google-user"
	}
	return result
}

func shortAccountSuffix(accountID string) string {
	clean := strings.TrimSpace(accountID)
	if len(clean) > 8 {
		clean = clean[:8]
	}
	clean = strings.ReplaceAll(clean, " ", "")
	return clean
}

func tokenValue(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprint(v)
}
