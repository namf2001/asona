package tasks

import "errors"

var (
	// ErrTaskNotFound is returned when a task is not found.
	ErrTaskNotFound = errors.New("task not found")
	
	// ErrPermissionDenied is returned when a user lacks needed permissions.
	ErrPermissionDenied = errors.New("permission denied for this task")
)
