package auth

import "errors"

var (
	errUnauthorized = errors.New("unauthorized")
	errUserNotFound = errors.New("user not found")
)
