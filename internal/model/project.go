package model

import "time"

// Project represents a project entity
type Project struct {
	ID          int64         `json:"id,omitempty"           db:"id"`
	WorkplaceID int64         `json:"workplace_id,omitempty" db:"workplace_id"`
	Name        string        `json:"name,omitempty"         db:"name"`
	Description string        `json:"description,omitempty"  db:"description"`
	Color       string        `json:"color,omitempty"        db:"color"`
	Access      ProjectAccess `json:"access,omitempty"       db:"access"`
	CreatedBy   int64         `json:"created_by,omitempty"   db:"created_by"`
	CreatedAt   time.Time     `json:"created_at,omitempty"   db:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at,omitempty"   db:"updated_at"`
}

// ProjectMember represents a user's membership in a project
type ProjectMember struct {
	ID        int64       `json:"id,omitempty"         db:"id"`
	ProjectID int64       `json:"project_id,omitempty" db:"project_id"`
	UserID    int64       `json:"user_id,omitempty"    db:"user_id"`
	Role      ProjectRole `json:"role,omitempty"       db:"role"`
	JoinedAt  time.Time   `json:"joined_at,omitempty"  db:"joined_at"`
}

// ProjectStatus represents a column/status within a project
type ProjectStatus struct {
	ID        int64     `json:"id,omitempty"         db:"id"`
	ProjectID int64     `json:"project_id,omitempty" db:"project_id"`
	Name      string    `json:"name,omitempty"       db:"name"`
	Color     string    `json:"color,omitempty"      db:"color"`
	Position  int       `json:"position,omitempty"   db:"position"`
	IsDefault bool      `json:"is_default,omitempty" db:"is_default"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
}
