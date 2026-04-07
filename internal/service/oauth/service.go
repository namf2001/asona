package oauth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"asona/config"
)

// Service defines basic OAuth capabilities.
type Service interface {
	GoogleConfig() *oauth2.Config
}

type service struct {
	googleConfig *oauth2.Config
}

func New() Service {
	cfg := config.GetConfig()
	googleConfig := &oauth2.Config{
		RedirectURL:  cfg.GoogleRedirectURL,
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &service{
		googleConfig: googleConfig,
	}
}

func (s *service) GoogleConfig() *oauth2.Config {
	return s.googleConfig
}
