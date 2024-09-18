package repository

import (
	"context"
	"github.com/wDRxxx/eventflow-backend/internal/models"
)

type Repository interface {
	EventByURLTitle(ctx context.Context, urlTitle string) (*models.Event, error)
	InsertEvent(ctx context.Context, event *models.Event) (int64, error)

	InsertUser(ctx context.Context, user *models.User) (int64, error)
}
