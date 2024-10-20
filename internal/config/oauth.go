package config

import (
	"os"
)

type OAuthConfig struct {
	callbackURL string
	redirectURL string

	googleClientID     string
	googleClientSecret string
}

func (c *OAuthConfig) RedirectURL() string {
	return c.redirectURL
}

func (c *OAuthConfig) GoogleClientID() string {
	return c.googleClientID
}

func (c *OAuthConfig) GoogleClientSecret() string {
	return c.googleClientSecret
}

func (c *OAuthConfig) CallbackURL() string {
	return c.callbackURL
}

func NewOAuthConfig() *OAuthConfig {
	callbackURL := os.Getenv("CALLBACK_URL")
	if callbackURL == "" {
		panic("CALLBACK_URL environment variable is empty")
	}

	redirectURL := os.Getenv("REDIRECT_URL")
	if redirectURL == "" {
		panic("REDIRECT_URL environment variable is empty")
	}

	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	if googleClientID == "" {
		panic("GOOGLE_CLIENT_ID environment variable is empty")
	}

	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if googleClientSecret == "" {
		panic("GOOGLE_CLIENT_SECRET environment variable is empty")
	}

	return &OAuthConfig{
		callbackURL:        callbackURL,
		redirectURL:        redirectURL,
		googleClientID:     googleClientID,
		googleClientSecret: googleClientSecret,
	}
}
