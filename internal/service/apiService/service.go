package apiService

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/wDRxxx/eventflow-backend/internal/config"
	"github.com/wDRxxx/eventflow-backend/internal/models"
	"github.com/wDRxxx/eventflow-backend/internal/repository"
	"github.com/wDRxxx/eventflow-backend/internal/service"
	"github.com/wDRxxx/eventflow-backend/internal/utils"
)

type serv struct {
	repo       repository.Repository
	authConfig *config.AuthConfig
}

func NewApiService(
	repo repository.Repository,
	authConfig *config.AuthConfig,
) service.ApiService {
	return &serv{
		repo:       repo,
		authConfig: authConfig,
	}
}

func (s *serv) Event(ctx context.Context, urlTitle string) (*models.Event, error) {
	event, err := s.repo.EventByURLTitle(ctx, urlTitle)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *serv) CreteEvent(ctx context.Context, event *models.Event) (int64, error) {
	if !event.IsPublic {
		event.URLTitle = uuid.NewString()
	}

	if !event.IsFree && len(event.Prices) == 0 {
		return 0, service.ErrNoPrices
	}

	if event.IsFree && len(event.Prices) > 0 {
		return 0, service.ErrPricesForFree
	}

	id, err := s.repo.InsertEvent(ctx, event)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *serv) UpdateEvent(ctx context.Context, event *models.Event) error {
	event.UpdatedAt = time.Now()
	err := s.repo.UpdateEvent(ctx, event)
	if err != nil {
		return err
	}

	return nil
}

func (s *serv) DeleteEvent(ctx context.Context, urlTitle string) error {
	err := s.repo.DeleteEvent(ctx, urlTitle)
	if err != nil {
		return err
	}

	return nil
}

//

func (s *serv) Ticket(ctx context.Context, ticketID string) (*models.Ticket, error) {
	ticket, err := s.repo.Ticket(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (s *serv) RegisterUser(ctx context.Context, user *models.User) error {
	_, err := s.repo.User(ctx, user.Email)
	if err == nil {
		return service.ErrUserAlreadyExists
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return err
	}
	user.Password = string(pass)

	_, err = s.repo.InsertUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *serv) Login(ctx context.Context, user *models.User) (string, error) {
	u, err := s.repo.User(ctx, user.Email)
	if err != nil {
		return "", service.ErrWrongCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		return "", service.ErrWrongCredentials
	}

	accessToken, err := utils.GenerateToken(
		u.Email,
		s.authConfig.RefreshTokenSecret,
		s.authConfig.RefreshTokenTTL,
	)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *serv) AccessToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := utils.VerifyToken(refreshToken, s.authConfig.RefreshTokenSecret)
	if err != nil {
		return "", err
	}

	accessToken, err := utils.GenerateToken(
		claims.Subject,
		s.authConfig.AccessTokenSecret,
		s.authConfig.AccessTokenTTL,
	)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
