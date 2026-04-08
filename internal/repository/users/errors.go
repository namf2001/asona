package users

import "errors"

var (
	// ErrUserNotFound is returned when a user record is not found in the database.
	ErrUserNotFound = errors.New("user not found")
	
	// ErrEmailAlreadyExists is returned when attempting to register an email that is already taken.
	ErrEmailAlreadyExists = errors.New("email already exists")
)
