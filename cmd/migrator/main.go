package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/joho/godotenv"

	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	m, err := migrate.New("file://"+migrationsPath, connectionURL)
	if err != nil {
		log.Fatalf("error creating migrator: %v", err)
	}

	if action != "up" {
		err = m.Down()
		if err != nil {
			log.Fatalf("error while migrating: %v", err)
		}

		log.Println("migrations were applied successfully!")
		return
	}

	err = m.Up()
	if err != nil {
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
