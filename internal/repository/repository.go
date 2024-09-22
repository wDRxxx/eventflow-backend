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
	Ticket(ctx context.Context, ticketID string) (*models.Ticket, error)

	InsertUser(ctx context.Context, user *models.User) (int64, error)
	User(ctx context.Context, userEmail string) (*models.User, error)
	InsertYookassaSettings(ctx context.Context, settings *models.YookassaSettings) (int64, error)
}
