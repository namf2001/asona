package auth

import (
	"context"

	"asona/internal/model"
	"asona/internal/repository"
	"asona/internal/service/mail"
	"asona/internal/service/oauth"
)

// Controller defines the business logic methods for authentication.
type Controller interface {
	// RegisterStep1SendOTP sends a one-time password to the given email address
	// as the first step of the 3-step registration flow.
	RegisterStep1SendOTP(ctx context.Context, email string) error
	// RegisterStep2VerifyOTP validates the OTP entered by the user against the stored token.
	// Returns an error if the code is incorrect or has expired.
	RegisterStep2VerifyOTP(ctx context.Context, email, otp string) error
	// RegisterStep3Complete finalizes registration by creating the user account and
	// issuing an initial session. Returns the created user and a JWT session token.
	RegisterStep3Complete(ctx context.Context, input RegisterInput) (model.User, string, error)
	// Login authenticates a user with email and password.
	// Returns the authenticated user and a JWT session token on success.
	Login(ctx context.Context, input LoginInput) (model.User, string, error)
	// GetProfile retrieves the full profile of the authenticated user by their ID.
	GetProfile(ctx context.Context, userID int64) (model.User, error)
	// Logout invalidates the session associated with the given JWT token.
	Logout(ctx context.Context, token string) error
	// GoogleAuthURL builds and returns the Google OAuth2 consent screen URL.
	// The state parameter is used to prevent CSRF attacks.
	GoogleAuthURL(ctx context.Context, state string) (string, error)
	// GoogleCallback completes the Google OAuth2 flow using the authorization code
	// returned by Google. Creates or links the account and issues a session token.
	GoogleCallback(ctx context.Context, code string) (model.User, string, error)
	// GetOnboardingState returns the onboarding status snapshot for the authenticated user.
	GetOnboardingState(ctx context.Context, userID int64) (OnboardingState, error)
	// CompleteOnboard marks the given user as having finished the onboarding flow.
	CompleteOnboard(ctx context.Context, userID int64) error
	// ValidateToken verifies the JWT signature and confirms the session is still active
	// in the database. Returns the userID and email embedded in the token claims.
	ValidateToken(ctx context.Context, token string) (userID int64, email string, err error)
}

type impl struct {
	repo  repository.Registry
	oauth oauth.Service
	mail  mail.Service
}

// New creates a new authentication controller instance.
func New(repo repository.Registry, oauthSvc oauth.Service, mailSvc mail.Service) Controller {
	return impl{repo: repo, oauth: oauthSvc, mail: mailSvc}

}
