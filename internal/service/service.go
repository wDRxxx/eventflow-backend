package service

import (
	"context"

	"github.com/wDRxxx/eventflow-backend/internal/models"
)

type ApiService interface {
	Event(ctx context.Context, urlTitle string) (*models.Event, error)
	CreteEvent(ctx context.Context, event *models.Event) (int64, error)
	DeleteEvent(ctx context.Context, urlTitle string) error
	UpdateEvent(ctx context.Context, event *models.Event) error

	Ticket(ctx context.Context, ticketID string) (*models.Ticket, error)

	RegisterUser(ctx context.Context, user *models.User) error
	Login(ctx context.Context, user *models.User) (string, error)
	AccessToken(ctx context.Context, refreshToken string) (string, error)
}
