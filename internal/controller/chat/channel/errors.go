package channel

import "errors"

var (
	// ErrChannelNotFound is returned when a channel is not found.
	ErrChannelNotFound = errors.New("channel not found")
	
	// ErrPermissionDenied is returned when a user lacks needed permissions.
	ErrPermissionDenied = errors.New("permission denied for this channel")
)
