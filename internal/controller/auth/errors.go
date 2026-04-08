package auth

import "errors"

var (
	// ErrUserNotFound is returned when the user credentials or record cannot be found.
	ErrUserNotFound = errors.New("authentication failed: user not found")

	// ErrInvalidPassword is returned when the password doesn't match.
	ErrInvalidPassword = errors.New("authentication failed: invalid password")

	// ErrUserAlreadyExists is returned when an email/username is already taken.
	ErrUserAlreadyExists = errors.New("registration failed: email already exists")

	// ErrSessionNotFound is returned when the target session does not exist.
	ErrSessionNotFound = errors.New("session not found")

	// ErrOAuthNotConfigured is returned when Google OAuth cannot be initialized.
	ErrOAuthNotConfigured = errors.New("google oauth not configured")

	// ErrOAuthExchangeFailed is returned when the authorization code cannot be exchanged.
	ErrOAuthExchangeFailed = errors.New("google oauth exchange failed")

	// ErrOAuthUserInfoFailed is returned when Google profile data cannot be fetched.
	ErrOAuthUserInfoFailed = errors.New("google oauth userinfo failed")

	// ErrOAuthEmailNotVerified is returned when the Google account email is not verified.
	ErrOAuthEmailNotVerified = errors.New("google oauth email not verified")
)
