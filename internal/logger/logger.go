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
	"github.com/wDRxxx/eventflow-backend/internal/logger/pretty"
)

func SetupLogger(envLevel string, logsPath string) {
	var logger *slog.Logger

	switch envLevel {
	case "dev":
		h := pretty.NewHandler(
			os.Stdout,
			&pretty.Options{
				Format: "text",
				Level:  "debug",
				Pretty: true,
			},
		)
		logger = slog.New(slogx.Accumulator(slogx.NewChain(h,
			slogm.StacktraceOnError(),
		)))

	case "prod":
		t := time.Now().Format("02-01-06")
		logFilePath := fmt.Sprintf("%s/%s.log", logsPath, t)

		f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			log.Fatalf("error creating new log file: %v", err)
		}
		closer.Add(2, func() error {
			slog.Info("closing log file")
			return f.Close()
		})

		h := slog.NewJSONHandler(f, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})

		logger = slog.New(slogx.Accumulator(slogx.NewChain(h,
			slogm.StacktraceOnError(),
		)))
	}

	slog.SetDefault(logger)
}
