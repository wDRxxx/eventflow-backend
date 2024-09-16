package config

import (
	"net"
	"os"
)

type HttpConfig struct {
	Host string
	Port string
}

func NewHttpConfig() *HttpConfig {
	host := os.Getenv("HTTP_HOST")
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		panic("HTTP_PORT environment variable is empty")
	}

	return &HttpConfig{
		Host: host,
		Port: port,
	}
}

func (c *HttpConfig) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}
