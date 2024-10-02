package service

import (
	"context"

	"github.com/wDRxxx/eventflow-backend/internal/models"
)

type ApiService interface {
	Event(ctx context.Context, urlTitle string) (*models.Event, error)
	Events(ctx context.Context, page int) ([]*models.Event, error)
	UserEvents(ctx context.Context, userID int) ([]*models.Event, error)
	CreteEvent(ctx context.Context, event *models.Event) (int64, error)
	DeleteEvent(ctx context.Context, userID int64, urlTitle string) error
	UpdateEvent(ctx context.Context, userID int64, event *models.Event) error

	Ticket(ctx context.Context, ticketID string) (*models.Ticket, error)
	UserTickets(ctx context.Context, userID int64) ([]*models.Ticket, error)
	BuyTicket(ctx context.Context, req *models.BuyTicketRequest) (string, error)

	RegisterUser(ctx context.Context, user *models.User) error
	Login(ctx context.Context, user *models.User) (string, error)
	AccessToken(ctx context.Context, refreshToken string) (string, error)
	User(ctx context.Context, userEmail string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
}
