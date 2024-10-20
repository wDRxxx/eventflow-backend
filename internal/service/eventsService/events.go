package eventsService

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/wDRxxx/eventflow-backend/internal/models"
	"github.com/wDRxxx/eventflow-backend/internal/repository"
	"github.com/wDRxxx/eventflow-backend/internal/service"
)

type eventsServ struct {
	repo repository.Repository
}

func NewEventsService(
	repo repository.Repository,
) service.EventsService {
	s := &eventsServ{
		repo: repo,
	}

	return s
}

func (s *eventsServ) Event(ctx context.Context, urlTitle string) (*models.Event, error) {
	event, err := s.repo.EventByURLTitle(ctx, urlTitle)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *eventsServ) Events(ctx context.Context, page int) ([]*models.Event, error) {
	events, err := s.repo.Events(ctx, page)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *eventsServ) CreateEvent(ctx context.Context, event *models.Event) (int64, error) {
	event.URLTitle = uuid.NewString()

	if !event.IsFree && len(event.Prices) == 0 {
		return 0, service.ErrNoPrices
	}

	if event.IsFree && len(event.Prices) > 0 {
		return 0, service.ErrPricesForFree
	}

	if event.Capacity == 0 {
		event.Capacity = 1000000000
	}

	id, err := s.repo.InsertEvent(ctx, event)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *eventsServ) UpdateEvent(ctx context.Context, userID int64, event *models.Event) error {
	e, err := s.repo.EventByURLTitle(ctx, event.URLTitle)
	if err != nil {
		return err
	}

	if e.CreatorID != userID {
		return service.ErrPermissionDenied
	}
	if len(event.Prices) == 0 {
		event.IsFree = true
	}

	event.UpdatedAt = time.Now()
	err = s.repo.UpdateEvent(ctx, event)
	if err != nil {
		return err
	}

	return nil
}

func (s *eventsServ) DeleteEvent(ctx context.Context, userID int64, urlTitle string) error {
	event, err := s.repo.EventByURLTitle(ctx, urlTitle)
	if err != nil {
		return err
	}

	if event.CreatorID != userID {
		return service.ErrPermissionDenied
	}

	err = s.repo.DeleteEvent(ctx, urlTitle)
	if err != nil {
		return err
	}

	return nil
}

func (s *eventsServ) UserEvents(ctx context.Context, userID int64) ([]*models.Event, error) {
	events, err := s.repo.UserEvents(ctx, userID)
	if err != nil {
		return nil, err
	}

	return events, nil
}
