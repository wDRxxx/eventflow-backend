package config

import (
	"os"
	"time"

	"github.com/xhit/go-str2duration/v2"
)

type AuthConfig struct {
	AccessTokenSecret  string
	AccessTokenTTL     time.Duration
	RefreshTokenSecret string
	RefreshTokenTTL    time.Duration
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

	return &AuthConfig{
		AccessTokenSecret:  ats,
		AccessTokenTTL:     attl,
		RefreshTokenSecret: rts,
		RefreshTokenTTL:    rttl,
	}
}
