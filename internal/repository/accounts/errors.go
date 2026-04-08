package accounts

import "errors"

var (
	// ErrAccountNotFound is returned when an OAuth provider link cannot be found.
	ErrAccountNotFound = errors.New("account not found")

	// ErrAccountAlreadyExists is returned when the OAuth provider link already exists.
	ErrAccountAlreadyExists = errors.New("account already exists")
)
