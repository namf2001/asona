package model

import "time"

// Channel represents an asynchronous communication room
type Channel struct {
	ID          int64       `json:"id,omitempty"           db:"id"`
	WorkplaceID int64       `json:"workplace_id,omitempty" db:"workplace_id"`
	ProjectID   int64       `json:"project_id,omitempty"   db:"project_id"`
	Name        string      `json:"name,omitempty"         db:"name"`
	Type        ChannelType `json:"type,omitempty"         db:"type"`
	CreatedBy   int64       `json:"created_by,omitempty"   db:"created_by"`
	CreatedAt   time.Time   `json:"created_at,omitempty"   db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at,omitempty"   db:"updated_at"`
}

// ChannelMember represents membership within a channel
type ChannelMember struct {
	ID         int64     `json:"id,omitempty"              db:"id"`
	ChannelID  int64     `json:"channel_id,omitempty"      db:"channel_id"`
	UserID     int64     `json:"user_id,omitempty"         db:"user_id"`
	LastReadAt time.Time `json:"last_read_at,omitempty"    db:"last_read_at"`
	JoinedAt   time.Time `json:"joined_at,omitempty"       db:"joined_at"`
}

// Message represents a message or a thread reply
type Message struct {
	ID        int64     `json:"id,omitempty"         db:"id"`
	ChannelID int64     `json:"channel_id,omitempty" db:"channel_id"`
	SenderID  int64     `json:"sender_id,omitempty"  db:"sender_id"`
	ParentID  int64     `json:"parent_id,omitempty"  db:"parent_id"` // Threading support
	Content   string    `json:"content,omitempty"    db:"content"`
	IsEdited  bool      `json:"is_edited,omitempty"  db:"is_edited"`
	IsDeleted bool      `json:"is_deleted,omitempty" db:"is_deleted"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}
