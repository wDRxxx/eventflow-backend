package config

import (
	"os"
	"time"

	"github.com/xhit/go-str2duration/v2"
)

type AuthConfig struct {
	accessTokenSecret  string
	accessTokenTTL     time.Duration
	refreshTokenSecret string
	refreshTokenTTL    time.Duration
	domain             string
}

func (c *AuthConfig) Domain() string {
	return c.domain
}

func (c *AuthConfig) AccessTokenSecret() string {
	return c.accessTokenSecret
}

func (c *AuthConfig) AccessTokenTTL() time.Duration {
	return c.accessTokenTTL
}

func (c *AuthConfig) RefreshTokenSecret() string {
	return c.refreshTokenSecret
}

func (c *AuthConfig) RefreshTokenTTL() time.Duration {
	return c.refreshTokenTTL
}

func NewAuthConfig() *AuthConfig {
	ats := os.Getenv("ACCESS_TOKEN_SECRET")
	if ats == "" {
		panic("ACCESS_TOKEN_SECRET environment variable is not set")
	}

	attl, err := str2duration.ParseDuration(os.Getenv("ACCESS_TOKEN_TTL"))
	if err != nil {
		panic("ACCESS_TOKEN_TTL environment variable is not set or has wrong format")
	}

	rts := os.Getenv("REFRESH_TOKEN_SECRET")
	if rts == "" {
		panic("REFRESH_TOKEN_SECRET environment variable is not set")
	}

	rttl, err := str2duration.ParseDuration(os.Getenv("REFRESH_TOKEN_TTL"))
	if err != nil {
		panic("REFRESH_TOKEN_TTL environment variable is not set or has wrong format")
	}

	domain := os.Getenv("DOMAIN")
	if domain == "" {
		panic("DOMAIN environment variable is not set")
	}

	return &AuthConfig{
		accessTokenSecret:  ats,
		accessTokenTTL:     attl,
		refreshTokenSecret: rts,
		refreshTokenTTL:    rttl,
		domain:             domain,
	}
}
