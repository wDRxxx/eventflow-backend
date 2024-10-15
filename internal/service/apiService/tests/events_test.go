package tests

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/wDRxxx/eventflow-backend/internal/closer"
	"github.com/wDRxxx/eventflow-backend/internal/config"
	"github.com/wDRxxx/eventflow-backend/internal/mailer/smtp"
	"github.com/wDRxxx/eventflow-backend/internal/models"
	"github.com/wDRxxx/eventflow-backend/internal/repository"
	"github.com/wDRxxx/eventflow-backend/internal/repository/mocks"
	"github.com/wDRxxx/eventflow-backend/internal/service"
	"github.com/wDRxxx/eventflow-backend/internal/service/apiService"
)

func TestEvent(t *testing.T) {
	t.Parallel()

	type repositoryMockFunc func(mc *minimock.Controller) repository.Repository

	var (
		wg  = &sync.WaitGroup{}
		ctx = context.Background()
		mc  = minimock.NewController(t)

		authCfg   = config.NewAuthConfig()
		mailerCfg = config.NewMailerConfig()
		mail, _   = smtp.NewSMTPMailer(mailerCfg, wg)

		repoErr  = errors.New("repo err")
		urlTitle = gofakeit.UUID()
		event    = &models.Event{
			ID:            gofakeit.Int64(),
			Title:         gofakeit.BeerName(),
			URLTitle:      urlTitle,
			Description:   gofakeit.ProductDescription(),
			BeginningTime: time.Now(),
			EndTime:       time.Now(),
			CreatorID:     gofakeit.Int64(),
			IsPublic:      gofakeit.Bool(),
			Location:      gofakeit.City(),
			IsFree:        gofakeit.Bool(),
			PreviewImage:  gofakeit.UUID(),
			UTCOffset:     gofakeit.Int64(),
			Capacity:      gofakeit.Int64(),
			MinimalAge:    gofakeit.Int64(),
			Prices:        nil,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
	)
	closer.SetGlobalCloser(closer.New(wg))

	tests := []struct {
		name           string
		want           *models.Event
		err            error
		repositoryMock repositoryMockFunc
	}{
		{
			name: "success case",
			want: event,
			err:  nil,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.EventByURLTitleMock.Expect(ctx, urlTitle).Return(event, nil)
				return mock
			},
		},
		{
			name: "failure case",
			want: nil,
			err:  repoErr,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.EventByURLTitleMock.Expect(ctx, urlTitle).Return(nil, repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryMock := tt.repositoryMock(mc)

			service := apiService.NewApiService(wg, repositoryMock, authCfg, mail)
			event, err := service.Event(ctx, urlTitle)

			require.Equal(t, tt.want, event)
			require.Equal(t, tt.err, err)
		})
	}
}

func TestEvents(t *testing.T) {
	t.Parallel()

	type repositoryMockFunc func(mc *minimock.Controller) repository.Repository

	var (
		wg  = &sync.WaitGroup{}
		ctx = context.Background()
		mc  = minimock.NewController(t)

		authCfg   = config.NewAuthConfig()
		mailerCfg = config.NewMailerConfig()
		mail, _   = smtp.NewSMTPMailer(mailerCfg, wg)

		repoErr = errors.New("repo err")
		page    = gofakeit.Int()
		events  = []*models.Event{{
			ID:            gofakeit.Int64(),
			Title:         gofakeit.BeerName(),
			URLTitle:      gofakeit.UUID(),
			Description:   gofakeit.ProductDescription(),
			BeginningTime: time.Now(),
			EndTime:       time.Now(),
			CreatorID:     gofakeit.Int64(),
			IsPublic:      gofakeit.Bool(),
			Location:      gofakeit.City(),
			IsFree:        gofakeit.Bool(),
			PreviewImage:  gofakeit.UUID(),
			UTCOffset:     gofakeit.Int64(),
			Capacity:      gofakeit.Int64(),
			MinimalAge:    gofakeit.Int64(),
			Prices:        nil,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}}
	)
	closer.SetGlobalCloser(closer.New(wg))

	tests := []struct {
		name           string
		want           []*models.Event
		err            error
		repositoryMock repositoryMockFunc
	}{
		{
			name: "success case",
			want: events,
			err:  nil,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.EventsMock.Expect(ctx, page).Return(events, nil)
				return mock
			},
		},
		{
			name: "failure case",
			want: nil,
			err:  repoErr,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.EventsMock.Expect(ctx, page).Return(nil, repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryMock := tt.repositoryMock(mc)

			service := apiService.NewApiService(wg, repositoryMock, authCfg, mail)
			event, err := service.Events(ctx, page)

			require.Equal(t, tt.want, event)
			require.Equal(t, tt.err, err)
		})
	}
}

func TestCreateEvent(t *testing.T) {
	t.Parallel()

	type repositoryMockFunc func(mc *minimock.Controller) repository.Repository

	var (
		wg  = &sync.WaitGroup{}
		ctx = context.Background()
		mc  = minimock.NewController(t)

		authCfg   = config.NewAuthConfig()
		mailerCfg = config.NewMailerConfig()
		mail, _   = smtp.NewSMTPMailer(mailerCfg, wg)

		repoErr = errors.New("repo err")

		id     = gofakeit.Int64()
		event1 = &models.Event{
			ID:            id,
			Title:         gofakeit.BeerName(),
			URLTitle:      gofakeit.UUID(),
			Description:   gofakeit.ProductDescription(),
			BeginningTime: time.Now(),
			EndTime:       time.Now(),
			CreatorID:     gofakeit.Int64(),
			IsPublic:      gofakeit.Bool(),
			Location:      gofakeit.City(),
			IsFree:        true,
			PreviewImage:  gofakeit.UUID(),
			UTCOffset:     gofakeit.Int64(),
			Capacity:      gofakeit.Int64(),
			MinimalAge:    gofakeit.Int64(),
			Prices:        nil,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		event2 = &models.Event{
			ID:            id,
			Title:         gofakeit.BeerName(),
			URLTitle:      gofakeit.UUID(),
			Description:   gofakeit.ProductDescription(),
			BeginningTime: time.Now(),
			EndTime:       time.Now(),
			CreatorID:     gofakeit.Int64(),
			IsPublic:      gofakeit.Bool(),
			Location:      gofakeit.City(),
			IsFree:        false,
			PreviewImage:  gofakeit.UUID(),
			UTCOffset:     gofakeit.Int64(),
			Capacity:      gofakeit.Int64(),
			MinimalAge:    gofakeit.Int64(),
			Prices:        nil,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		event3 = &models.Event{
			ID:            id,
			Title:         gofakeit.BeerName(),
			URLTitle:      gofakeit.UUID(),
			Description:   gofakeit.ProductDescription(),
			BeginningTime: time.Now(),
			EndTime:       time.Now(),
			CreatorID:     gofakeit.Int64(),
			IsPublic:      gofakeit.Bool(),
			Location:      gofakeit.City(),
			IsFree:        true,
			PreviewImage:  gofakeit.UUID(),
			UTCOffset:     gofakeit.Int64(),
			Capacity:      gofakeit.Int64(),
			MinimalAge:    gofakeit.Int64(),
			Prices: []*models.Price{{
				ID:       id,
				Price:    gofakeit.Int64(),
				Currency: gofakeit.CurrencyShort(),
			}},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	)
	closer.SetGlobalCloser(closer.New(wg))

	tests := []struct {
		name           string
		want           int64
		event          *models.Event
		err            error
		repositoryMock repositoryMockFunc
	}{
		{
			name:  "success case",
			want:  id,
			event: event1,
			err:   nil,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.InsertEventMock.Expect(ctx, event1).Return(id, nil)
				return mock
			},
		},
		{
			name:  "wrong prices case 1",
			want:  0,
			event: event2,
			err:   service.ErrNoPrices,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				return mock
			},
		},
		{
			name:  "wrong prices case 2",
			want:  0,
			event: event3,
			err:   service.ErrPricesForFree,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				return mock
			},
		},
		{
			name:  "failure case",
			want:  0,
			err:   repoErr,
			event: event1,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.InsertEventMock.Expect(ctx, event1).Return(id, repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryMock := tt.repositoryMock(mc)

			service := apiService.NewApiService(wg, repositoryMock, authCfg, mail)
			eventID, err := service.CreateEvent(ctx, tt.event)

			require.Equal(t, tt.want, eventID)
			require.Equal(t, tt.err, err)
		})
	}
}

func TestUpdateEvent(t *testing.T) {
	t.Parallel()

	type repositoryMockFunc func(mc *minimock.Controller) repository.Repository

	var (
		wg  = &sync.WaitGroup{}
		ctx = context.Background()
		mc  = minimock.NewController(t)

		authCfg   = config.NewAuthConfig()
		mailerCfg = config.NewMailerConfig()
		mail, _   = smtp.NewSMTPMailer(mailerCfg, wg)

		repoErr = errors.New("repo err")

		id     = gofakeit.Int64()
		userID = gofakeit.Int64()
		event1 = &models.Event{
			ID:            id,
			Title:         gofakeit.BeerName(),
			URLTitle:      gofakeit.UUID(),
			Description:   gofakeit.ProductDescription(),
			BeginningTime: time.Now(),
			EndTime:       time.Now(),
			CreatorID:     userID,
			IsPublic:      gofakeit.Bool(),
			Location:      gofakeit.City(),
			IsFree:        true,
			PreviewImage:  gofakeit.UUID(),
			UTCOffset:     gofakeit.Int64(),
			Capacity:      gofakeit.Int64(),
			MinimalAge:    gofakeit.Int64(),
			Prices:        nil,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		event2 = &models.Event{
			ID:            id,
			Title:         gofakeit.BeerName(),
			URLTitle:      gofakeit.UUID(),
			Description:   gofakeit.ProductDescription(),
			BeginningTime: time.Now(),
			EndTime:       time.Now(),
			CreatorID:     userID,
			IsPublic:      gofakeit.Bool(),
			Location:      gofakeit.City(),
			IsFree:        false,
			PreviewImage:  gofakeit.UUID(),
			UTCOffset:     gofakeit.Int64(),
			Capacity:      gofakeit.Int64(),
			MinimalAge:    gofakeit.Int64(),
			Prices: []*models.Price{{
				ID:       id,
				Price:    gofakeit.Int64(),
				Currency: gofakeit.CurrencyShort(),
			}},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		event3 = &models.Event{
			ID:            id,
			Title:         gofakeit.BeerName(),
			URLTitle:      gofakeit.UUID(),
			Description:   gofakeit.ProductDescription(),
			BeginningTime: time.Now(),
			EndTime:       time.Now(),
			CreatorID:     gofakeit.Int64(),
			IsPublic:      gofakeit.Bool(),
			Location:      gofakeit.City(),
			IsFree:        false,
			PreviewImage:  gofakeit.UUID(),
			UTCOffset:     gofakeit.Int64(),
			Capacity:      gofakeit.Int64(),
			MinimalAge:    gofakeit.Int64(),
			Prices: []*models.Price{{
				ID:       id,
				Price:    gofakeit.Int64(),
				Currency: gofakeit.CurrencyShort(),
			}},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	)
	closer.SetGlobalCloser(closer.New(wg))

	tests := []struct {
		name           string
		event          *models.Event
		err            error
		repositoryMock repositoryMockFunc
	}{
		{
			name:  "success case",
			err:   nil,
			event: event1,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.EventByURLTitleMock.Expect(ctx, event1.URLTitle).Return(event1, nil)
				mock.UpdateEventMock.Expect(ctx, event1).Return(nil)
				return mock
			},
		},
		{
			name:  "wrong prices case 1",
			err:   nil,
			event: event2,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.EventByURLTitleMock.Expect(ctx, event2.URLTitle).Return(event2, nil)
				mock.UpdateEventMock.Expect(ctx, event2).Return(nil)

				return mock
			},
		},
		{
			name:  "wrong user case",
			err:   service.ErrPermissionDenied,
			event: event3,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.EventByURLTitleMock.Expect(ctx, event3.URLTitle).Return(event3, nil)
				return mock
			},
		},
		{
			name:  "failure case",
			err:   repoErr,
			event: event2,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.EventByURLTitleMock.Expect(ctx, event2.URLTitle).Return(event2, nil)
				mock.UpdateEventMock.Expect(ctx, event2).Return(repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryMock := tt.repositoryMock(mc)

			service := apiService.NewApiService(wg, repositoryMock, authCfg, mail)
			err := service.UpdateEvent(ctx, userID, tt.event)

			require.Equal(t, tt.err, err)
		})
	}
}

func TestDeleteEvent(t *testing.T) {
	t.Parallel()

	type repositoryMockFunc func(mc *minimock.Controller) repository.Repository

	var (
		wg  = &sync.WaitGroup{}
		ctx = context.Background()
		mc  = minimock.NewController(t)

		authCfg   = config.NewAuthConfig()
		mailerCfg = config.NewMailerConfig()
		mail, _   = smtp.NewSMTPMailer(mailerCfg, wg)

		repoErr = errors.New("repo err")

		urlTitle = gofakeit.UUID()
		userID   = gofakeit.Int64()
		event    = &models.Event{
			ID:            gofakeit.Int64(),
			Title:         gofakeit.BeerName(),
			URLTitle:      urlTitle,
			Description:   gofakeit.ProductDescription(),
			BeginningTime: time.Now(),
			EndTime:       time.Now(),
			CreatorID:     userID,
			IsPublic:      gofakeit.Bool(),
			Location:      gofakeit.City(),
			IsFree:        true,
			PreviewImage:  gofakeit.UUID(),
			UTCOffset:     gofakeit.Int64(),
			Capacity:      gofakeit.Int64(),
			MinimalAge:    gofakeit.Int64(),
			Prices:        nil,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
	)
	closer.SetGlobalCloser(closer.New(wg))

	tests := []struct {
		name   string
		err    error
		userID int64

		repositoryMock repositoryMockFunc
	}{
		{
			name:   "success case",
			err:    nil,
			userID: userID,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.EventByURLTitleMock.Expect(ctx, urlTitle).Return(event, nil)
				mock.DeleteEventMock.Expect(ctx, urlTitle).Return(nil)
				return mock
			},
		},
		{
			name:   "wrong user case",
			err:    service.ErrPermissionDenied,
			userID: gofakeit.Int64(),
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.EventByURLTitleMock.Expect(ctx, urlTitle).Return(event, nil)
				return mock
			},
		},
		{
			name:   "failure case",
			err:    repoErr,
			userID: userID,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.EventByURLTitleMock.Expect(ctx, urlTitle).Return(event, nil)
				mock.DeleteEventMock.Expect(ctx, urlTitle).Return(repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryMock := tt.repositoryMock(mc)

			service := apiService.NewApiService(wg, repositoryMock, authCfg, mail)
			err := service.DeleteEvent(ctx, tt.userID, urlTitle)

			require.Equal(t, tt.err, err)
		})
	}
}

func TestUserEvents(t *testing.T) {
	t.Parallel()

	type repositoryMockFunc func(mc *minimock.Controller) repository.Repository

	var (
		wg  = &sync.WaitGroup{}
		ctx = context.Background()
		mc  = minimock.NewController(t)

		authCfg   = config.NewAuthConfig()
		mailerCfg = config.NewMailerConfig()
		mail, _   = smtp.NewSMTPMailer(mailerCfg, wg)

		repoErr = errors.New("repo err")
		userID  = gofakeit.Int64()
		events  = []*models.Event{{
			ID:            gofakeit.Int64(),
			Title:         gofakeit.BeerName(),
			URLTitle:      gofakeit.UUID(),
			Description:   gofakeit.ProductDescription(),
			BeginningTime: time.Now(),
			EndTime:       time.Now(),
			CreatorID:     userID,
			IsPublic:      gofakeit.Bool(),
			Location:      gofakeit.City(),
			IsFree:        gofakeit.Bool(),
			PreviewImage:  gofakeit.UUID(),
			UTCOffset:     gofakeit.Int64(),
			Capacity:      gofakeit.Int64(),
			MinimalAge:    gofakeit.Int64(),
			Prices:        nil,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}}
	)
	closer.SetGlobalCloser(closer.New(wg))

	tests := []struct {
		name           string
		want           []*models.Event
		err            error
		repositoryMock repositoryMockFunc
	}{
		{
			name: "success case",
			want: events,
			err:  nil,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.UserEventsMock.Expect(ctx, userID).Return(events, nil)
				return mock
			},
		},
		{
			name: "failure case",
			want: nil,
			err:  repoErr,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.UserEventsMock.Expect(ctx, userID).Return(nil, repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryMock := tt.repositoryMock(mc)

			service := apiService.NewApiService(wg, repositoryMock, authCfg, mail)
			event, err := service.UserEvents(ctx, userID)

			require.Equal(t, tt.want, event)
			require.Equal(t, tt.err, err)
		})
	}
}
