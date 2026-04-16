package model

import (
	"time"

	"asona/internal/model/common"
)

// User represents a user in the system
type User struct {
	ID               int64            `json:"id,omitempty"              db:"id"`
	Name             string           `json:"name,omitempty"            db:"name"`
	Username         string           `json:"username,omitempty"        db:"username"`
	DisplayName      string           `json:"display_name,omitempty"    db:"display_name"`
	Email            string           `json:"email,omitempty"           db:"email"`
	EmailVerified    *time.Time       `json:"email_verified,omitempty"  db:"email_verified"`
	Image            string           `json:"image,omitempty"           db:"avatar_url"`
	Password         string           `json:"-"                         db:"password"`
	IsActive         bool             `json:"is_active,omitempty"       db:"is_active"`
	OnboardingStatus OnboardingStatus `json:"onboarding_status,omitempty" db:"onboarding_status"`
	OnboardingStep   int16            `json:"onboarding_step,omitempty"   db:"onboarding_step"`
	OnboardedAt      *time.Time       `json:"onboarded_at,omitempty"    db:"onboarded_at"`
	CreatedAt        time.Time        `json:"created_at,omitempty"      db:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at,omitempty"      db:"updated_at"`
}

// Prepare auto-sets timestamps before insert/update
func (u *User) Prepare() {
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	u.UpdatedAt = time.Now()
}

// Validate validates user data
func (u *User) Validate() error {
	if u.Email == "" {
		return common.ErrInvalidEmail
	}
	if u.Name == "" {
		return common.ErrInvalidName
	}
	return nil
}
