package tests

import (
	"github.com/wDRxxx/eventflow-backend/internal/config"
)

func init() {
	config.Load(".env")
}
