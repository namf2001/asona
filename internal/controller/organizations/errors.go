package organizations

import "errors"

var (
	// ErrOrganizationNotFound is returned when an organization does not exist.
	ErrOrganizationNotFound = errors.New("organization not found")
	
	// ErrMemberNotFound is returned when a user is not a member of an organization.
	ErrMemberNotFound = errors.New("organization member not found")
	
	// ErrPermissionDenied is returned when a user lacks permission for an action.
	ErrPermissionDenied = errors.New("permission denied for this organization")
)
