package model

import "time"

// Workplace represents an actual work environment within an organization
type Workplace struct {
	ID        int64         `json:"id,omitempty"          db:"id"`
	Name      string        `json:"name,omitempty"        db:"name"`
	IconURL   string        `json:"icon_url,omitempty"    db:"icon_url"`
	Size      WorkplaceSize `json:"size,omitempty"        db:"size"`
	CreatedBy int64         `json:"created_by,omitempty"  db:"created_by"`
	CreatedAt time.Time     `json:"created_at,omitempty"  db:"created_at"`
	UpdatedAt time.Time     `json:"updated_at,omitempty"  db:"updated_at"`
}

// WorkplaceMember represents a user's membership in a workplace
type WorkplaceMember struct {
	ID          int64         `json:"id,omitempty"              db:"id"`
	WorkplaceID int64         `json:"workplace_id,omitempty"    db:"workplace_id"`
	UserID      int64         `json:"user_id,omitempty"         db:"user_id"`
	Role        WorkplaceRole `json:"role,omitempty"            db:"role"`
	JoinedAt    time.Time     `json:"joined_at,omitempty"       db:"joined_at"`
}
