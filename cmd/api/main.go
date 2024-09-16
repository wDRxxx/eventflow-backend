package main

import (
	"context"
	"flag"
	"github.com/wDRxxx/eventflow-backend/internal/app"
	"log"
)

var envPath string

func init() {
	flag.StringVar(&envPath, "env-path", ".env", "path to .env file")

	flag.Parse()
}

func main() {
	ctx := context.Background()
	a, err := app.NewApp(ctx, envPath)
	if err != nil {
		log.Fatalf("error creating app: %v", err)
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("error running app: %v", err)
	}
}
