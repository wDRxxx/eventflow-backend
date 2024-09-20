package repository

import (
	"context"

	"github.com/wDRxxx/eventflow-backend/internal/models"
)

type Repository interface {
	EventByURLTitle(ctx context.Context, urlTitle string) (*models.Event, error)
	InsertEvent(ctx context.Context, event *models.Event) (int64, error)
	UpdateEvent(ctx context.Context, event *models.Event) error
	DeleteEvent(ctx context.Context, urlTitle string) error

	InsertTicket(ctx context.Context, ticket *models.Ticket) (string, error)

	InsertUser(ctx context.Context, user *models.User) (int64, error)
	InsertYookassaSettings(ctx context.Context, settings *models.YookassaSettings) (int64, error)
}
