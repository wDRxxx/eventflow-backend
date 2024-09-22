package main

import (
	"context"
	"flag"
	"log"

	"github.com/wDRxxx/eventflow-backend/internal/app"
	"github.com/wDRxxx/eventflow-backend/internal/logger"
)

var envPath, envLevel, logsPath string

func init() {
	flag.StringVar(&envPath, "env-path", ".env", "path to .env file")
	flag.StringVar(&envLevel, "env-level", "dev", "dev/prod")
	flag.StringVar(&logsPath, "logs-path", "./logs", "path to folder with logs")

	flag.Parse()
}

func main() {
	ctx := context.Background()
	logger.SetupLogger(envLevel, logsPath)

	a, err := app.NewApp(ctx, envPath)
	if err != nil {
		log.Fatalf("error creating app: %v", err)
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("error running app: %v", err)
	}
}
