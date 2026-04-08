package common

import "errors"

var (
	ErrInvalidEmail      = errors.New("invalid email")
	ErrInvalidName       = errors.New("invalid name")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserNotFound      = errors.New("user not found")
	ErrWorkspaceNotFound = errors.New("workspace not found")
)
