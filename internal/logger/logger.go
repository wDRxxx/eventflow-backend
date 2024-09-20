package logger

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/cappuccinotm/slogx"
	"github.com/cappuccinotm/slogx/slogm"

	"github.com/wDRxxx/eventflow-backend/internal/closer"
)

func SetupLogger(envLevel string, logsPath string) {
	var logger *slog.Logger

	switch envLevel {
	case "dev":
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		t := time.Now().Format("02-01-06T15-04-05")
		logFilePath := fmt.Sprintf("%s/%s.log", logsPath, t)

		f, err := os.Create(logFilePath)
		if err != nil {
			log.Fatalf("error creating new log file: %v", err)
		}
		closer.Add(f.Close)

		h := slog.NewJSONHandler(f, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})

		logger = slog.New(slogx.Accumulator(slogx.NewChain(h,
			slogm.StacktraceOnError(),
		)))
	}

	slog.SetDefault(logger)
}
