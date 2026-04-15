package workplaces

import "errors"

var (
	ErrWorkplaceNotFound = errors.New("workplace not found")
	ErrUnauthorized      = errors.New("unauthorized")
)
