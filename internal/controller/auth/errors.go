package auth

import "errors"

var (
	// ErrUserNotFound is returned when the user credentials or record cannot be found.
	ErrUserNotFound = errors.New("authentication failed: user not found")

	// ErrInvalidPassword is returned when the password doesn't match.
	ErrInvalidPassword = errors.New("authentication failed: invalid password")

	// ErrInvalidToken is returned when the JWT token is malformed, expired, or has an invalid signature.
	ErrInvalidToken = errors.New("authentication failed: invalid or expired token")

	// ErrUserAlreadyExists is returned when an email is already registered.
	ErrUserAlreadyExists = errors.New("registration failed: email already exists")

	// ErrUsernameAlreadyExists is returned when the chosen username is already taken.
	ErrUsernameAlreadyExists = errors.New("registration failed: username already exists")

	// ErrWrongOTP is returned when the OTP provided is invalid or expired.
	ErrWrongOTP = errors.New("registration failed: wrong or expired otp")

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
