package config

import (
	"log"
	"os"
)

type MailerConfig struct {
	login    string
	password string
	host     string
	port     string

	orderTemplatePath string
}

func (c *MailerConfig) OrderTemplatePath() string {
	return c.orderTemplatePath
}

func (c *MailerConfig) Login() string {
	return c.login
}

func (c *MailerConfig) Password() string {
	return c.password
}

func (c *MailerConfig) Host() string {
	return c.host
}

func (c *MailerConfig) Port() string {
	return c.port
}

func NewMailerConfig() *MailerConfig {
	login := os.Getenv("MAILER_LOGIN")
	if login == "" {
		log.Fatal("MAILER_LOGIN environment variable is empty")
	}
	password := os.Getenv("MAILER_PASSWORD")
	if password == "" {
		log.Fatal("MAILER_PASSWORD environment variable is empty")
	}
	host := os.Getenv("MAILER_HOST")
	if host == "" {
		log.Fatal("MAILER_HOST environment variable is empty")
	}
	port := os.Getenv("MAILER_PORT")
	if port == "" {
		log.Fatal("MAILER_PORT environment variable is empty")
	}
	orderTemplatePath := os.Getenv("MAILER_ORDER_TEMPLATE_PATH")
	if orderTemplatePath == "" {
		log.Fatal("MAILER_ORDER_TEMPLATE_PATH environment variable is empty")
	}

	return &MailerConfig{
		login:             login,
		password:          password,
		host:              host,
		port:              port,
		orderTemplatePath: orderTemplatePath,
	}
}
