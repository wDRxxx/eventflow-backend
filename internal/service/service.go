package service

import (
	"context"

	"github.com/wDRxxx/eventflow-backend/internal/models"
)

type ApiService interface {
	Event(ctx context.Context, urlTitle string) (*models.Event, error)
	CreteEvent(ctx context.Context, event *models.Event) (int64, error)
	DeleteEvent(ctx context.Context, urlTitle string) error
}
