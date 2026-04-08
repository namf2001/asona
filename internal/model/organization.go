package model

import "time"

// Organization represents a top-level organization entity
type Organization struct {
	ID          int64     `json:"id,omitempty"          db:"id"`
	Name        string    `json:"name,omitempty"        db:"name"`
	Slug        string    `json:"slug,omitempty"        db:"slug"`
	LogoURL     string    `json:"logo_url,omitempty"    db:"logo_url"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatedBy   int64     `json:"created_by,omitempty"  db:"created_by"`
	CreatedAt   time.Time `json:"created_at,omitempty"  db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"  db:"updated_at"`
}

// OrganizationMember represents a user's membership in an organization
type OrganizationMember struct {
	ID             int64     `json:"id,omitempty"              db:"id"`
	OrganizationID int64     `json:"organization_id,omitempty" db:"organization_id"`
	UserID         int64     `json:"user_id,omitempty"         db:"user_id"`
	Role           OrgRole   `json:"role,omitempty"            db:"role"`
	JoinedAt       time.Time `json:"joined_at,omitempty"       db:"joined_at"`
}
