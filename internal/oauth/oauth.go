package oauth

import (
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/wDRxxx/eventflow-backend/internal/config"
)

type OAuth struct {
	redirectURL string

	googleConfig *oauth2.Config
}

func (o *OAuth) RedirectURL() string {
	return o.redirectURL
}

func (o *OAuth) GoogleConfig() *oauth2.Config {
	return o.googleConfig
}

func NewOAuth(config *config.OAuthConfig) *OAuth {
	return &OAuth{
		redirectURL: config.RedirectURL(),
		googleConfig: &oauth2.Config{
			RedirectURL:  fmt.Sprintf(config.CallbackURL(), "google"),
			ClientID:     config.GoogleClientID(),
			ClientSecret: config.GoogleClientSecret(),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
			Endpoint:     google.Endpoint,
		},
	}
}
