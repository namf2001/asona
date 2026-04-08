package message

import "errors"

var (
	// ErrMessageNotFound is returned when a message is not found.
	ErrMessageNotFound = errors.New("message not found")
)
