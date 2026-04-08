package projects

import "errors"

var (
	// ErrProjectNotFound is returned when a project does not exist.
	ErrProjectNotFound = errors.New("project not found")
	
	// ErrProjectAlreadyExists is returned when a project with the same name exists.
	ErrProjectAlreadyExists = errors.New("project already exists")
)
