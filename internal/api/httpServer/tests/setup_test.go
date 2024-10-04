package tests

import (
	"github.com/wDRxxx/eventflow-backend/internal/config"
	"github.com/wDRxxx/eventflow-backend/internal/metrics"
)

func init() {
	metrics.Init("test")
	config.Load(".env")
}
