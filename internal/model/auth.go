package model

import "time"

// Account represents an OAuth provider account for a user
type Account struct {
	ID                int64      `json:"id,omitempty"                 db:"id"`
	UserID            int64      `json:"user_id,omitempty"            db:"user_id"`
	Provider          string     `json:"provider,omitempty"           db:"provider"`
	ProviderAccountID string     `json:"provider_account_id,omitempty" db:"provider_account_id"`
	AccessToken       string     `json:"access_token,omitempty"       db:"access_token"`
	RefreshToken      string     `json:"refresh_token,omitempty"      db:"refresh_token"`
	TokenExpiresAt    *time.Time  `json:"token_expires_at,omitempty"    db:"token_expires_at"`
	IDToken           string     `json:"id_token,omitempty"           db:"id_token"`
	Scope             string     `json:"scope,omitempty"              db:"scope"`
	CreatedAt         time.Time  `json:"created_at,omitempty"         db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at,omitempty"         db:"updated_at"`
}

// Session represents an active session in the database
type Session struct {
	ID           int64      `json:"id,omitempty"              db:"id"`
	UserID       int64      `json:"user_id,omitempty"         db:"user_id"`
	SessionToken string     `json:"session_token,omitempty"   db:"session_token"`
	ExpiresAt    time.Time  `json:"expires_at,omitempty"      db:"expires_at"`
	UserAgent    string     `json:"user_agent,omitempty"      db:"user_agent"`
	IPAddress    string     `json:"ip_address,omitempty"      db:"ip_address"`
	CreatedAt    time.Time  `json:"created_at,omitempty"      db:"created_at"`
}

// VerificationToken represents a short-lived token used for verification purposes
type VerificationToken struct {
	Identifier string    `json:"identifier" db:"identifier"`
	Token      string    `json:"token"      db:"token"`
	Expires    time.Time `json:"expires"    db:"expires"`
}
