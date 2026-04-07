package validation

import (
	"regexp"
	"strings"
)

// Email validation regex pattern
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Password validation regex pattern - minimum 6 characters
var passwordRegex = regexp.MustCompile(`^.{6,}$`)

// IsValidEmail validates email format using regex
func IsValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}
	return emailRegex.MatchString(email)
}

// IsValidPassword validates password format using regex - minimum 6 characters
func IsValidPassword(password string) bool {
	password = strings.TrimSpace(password)
	if password == "" {
		return false
	}
	return passwordRegex.MatchString(password)
}
