package tasks

import "errors"

var (
	// ErrTaskNotFound is returned when a task does not exist.
	ErrTaskNotFound = errors.New("task not found")
)
