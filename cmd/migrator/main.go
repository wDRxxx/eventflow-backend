package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func main() {
	var action, migrationsPath, envPath string

	flag.StringVar(&action, "action", "up", "up/down")
	flag.StringVar(&envPath, "env-path", "", "path to .env file")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to folder with migrations")
	flag.Parse()
	if migrationsPath == "" {
		panic("path to migrations folder is empty")
	}
	if envPath == "" {
		panic("path to .env file is empty")
	}

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	connectionURL := mustConnectionURL()

	var m *migrate.Migrate
	for range 10 {
		time.Sleep(1 * time.Second)
		m, err = migrate.New("file://"+migrationsPath, connectionURL)
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Fatalf("error creating migrator: %v", err)
	}

	if action != "up" {
		err = m.Down()

		log.Println("migrations were applied successfully!")
		return
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("error while migrating: %v", err)
	}
	log.Println("migrations were applied successfully!")
}

func mustConnectionURL() string {
	pgUser := os.Getenv("POSTGRES_USER")
	pgPass := os.Getenv("POSTGRES_PASSWORD")
	pgDB := os.Getenv("POSTGRES_DB")
	pgHost := os.Getenv("POSTGRES_HOST")
	pgPort := os.Getenv("POSTGRES_PORT")
	if pgUser == "" {
		panic("POSTGRES_USER environment variable is empty")
	}
	if pgPass == "" {
		panic("POSTGRES_PASSWORD environment variable is empty")
	}
	if pgDB == "" {
		panic("POSTGRES_DB environment variable is empty")
	}
	if pgHost == "" {
		panic("POSTGRES_HOST environment variable is empty")
	}
	if pgPort == "" {
		panic("POSTGRES_PORT environment variable is empty")
	}

	sourceURL := fmt.Sprintf(
		"pgx5://%s:%s@%s:%s/%s",
		pgUser,
		pgPass,
		pgHost,
		pgPort,
		pgDB,
	)

	return sourceURL
}
