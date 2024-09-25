package config

import (
	"net"
	"os"
	"strings"
)

type HttpConfig struct {
	Host    string
	Port    string
	Origins []string
}

func NewHttpConfig() *HttpConfig {
	host := os.Getenv("HTTP_HOST")
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		panic("HTTP_PORT environment variable is empty")
	}

	origins := os.Getenv("HTTP_ORIGINS")
	if len(origins) == 0 {
		panic("http origins not found")
	}

	return &HttpConfig{
		Host:    host,
		Port:    port,
		Origins: strings.Split(origins, " "),
	}
}

func (c *HttpConfig) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}
