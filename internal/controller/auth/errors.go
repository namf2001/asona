package auth

import "errors"

var (
	// ErrUserNotFound is returned when the user credentials or record cannot be found.
	ErrUserNotFound = errors.New("authentication failed: user not found")
	
	// ErrInvalidPassword is returned when the password doesn't match.
	ErrInvalidPassword = errors.New("authentication failed: invalid password")
	
	// ErrUserAlreadyExists is returned when an email/username is already taken.
	ErrUserAlreadyExists = errors.New("registration failed: email already exists")
)
