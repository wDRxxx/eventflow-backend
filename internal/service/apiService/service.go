package apiService

import (
	"context"

	"github.com/wDRxxx/eventflow-backend/internal/models"
	"github.com/wDRxxx/eventflow-backend/internal/repository"
	"github.com/wDRxxx/eventflow-backend/internal/service"
)

type serv struct {
	repo repository.Repository
}

func NewApiService(repo repository.Repository) service.ApiService {
	return &serv{
		repo: repo,
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
	id, err := s.repo.InsertEvent(ctx, event)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *serv) UpdateEvent(ctx context.Context, event *models.Event) error {
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
