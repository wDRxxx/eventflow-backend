package config

import (
	"net"
	"os"
	"strings"
)

type HttpConfig struct {
	host      string
	port      string
	origins   []string
	staticDir string
}

func (c *HttpConfig) StaticDir() string {
	return c.staticDir
}

func (c *HttpConfig) Origins() []string {
	return c.origins
}

func NewHttpConfig() *HttpConfig {
	host := os.Getenv("HTTP_HOST")
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		panic("HTTP_PORT environment variable is empty")
	}

	staticDir := os.Getenv("HTTP_STATIC_DIR")
	if staticDir == "" {
		panic("HTTP_STATIC_DIR environment variable is empty")
	}

	origins := os.Getenv("HTTP_ORIGINS")
	if len(origins) == 0 {
		panic("HTTP_ORIGINS environment variable is empty")
	}

	return &HttpConfig{
		host:      host,
		port:      port,
		origins:   strings.Split(origins, " "),
		staticDir: staticDir,
	}
}

func (c *HttpConfig) Address() string {
	return net.JoinHostPort(c.host, c.port)
}
