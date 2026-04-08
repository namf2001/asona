package organizations

import "errors"

var (
	// ErrOrganizationNotFound is returned when an organization record is not found.
	ErrOrganizationNotFound = errors.New("organization not found")
	
	// ErrMemberNotFound is returned when an organization member record is not found.
	ErrMemberNotFound = errors.New("member not found")
)
