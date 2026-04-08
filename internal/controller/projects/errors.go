package projects

import "errors"

var (
	// ErrProjectNotFound is returned when a project is not found.
	ErrProjectNotFound = errors.New("project not found")
	
	// ErrPermissionDenied is returned when a user lacks needed permissions.
	ErrPermissionDenied = errors.New("permission denied for this project")
)
