package config

import (
	"fmt"
	"os"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func NewPostgresConfig() *PostgresConfig {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		panic("POSTGRES_HOST environment variable is empty")
	}

	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		panic("POSTGRES_PORT environment variable is empty")
	}

	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		panic("POSTGRES_USER environment variable is empty")
	}

	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		panic("POSTGRES_PASSWORD environment variable is empty")
	}

	database := os.Getenv("POSTGRES_DB")
	if database == "" {
		panic("POSTGRES_DB environment variable is empty")
	}

	return &PostgresConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Database: database,
	}
}

func (s *PostgresConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		s.Host,
		s.Port,
		s.Database,
		s.User,
		s.Password,
	)
}
