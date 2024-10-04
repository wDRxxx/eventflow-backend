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
	"github.com/wDRxxx/yookassa-go-sdk/yookassa"

	"github.com/wDRxxx/eventflow-backend/internal/closer"
	"github.com/wDRxxx/eventflow-backend/internal/config"
	"github.com/wDRxxx/eventflow-backend/internal/mailer/smtp"
	"github.com/wDRxxx/eventflow-backend/internal/models"
	"github.com/wDRxxx/eventflow-backend/internal/repository"
	"github.com/wDRxxx/eventflow-backend/internal/repository/mocks"
	"github.com/wDRxxx/eventflow-backend/internal/service/apiService"
)

func _TestBuyTicket(t *testing.T) {
	t.Parallel()

	type repositoryMockFunc func(mc *minimock.Controller) repository.Repository

	var (
		wg  = &sync.WaitGroup{}
		ctx = context.Background()
		mc  = minimock.NewController(t)

		authCfg   = config.NewAuthConfig()
		yooCfg    = config.NewYookassaConfig()
		mailerCfg = config.NewMailerConfig()
		yooClient = yookassa.NewClient(yooCfg.ShopID(), yooCfg.ShopKey())
		mail, _   = smtp.NewSMTPMailer(mailerCfg, wg)

		ticketID = gofakeit.UUID()
		//repoErr  = errors.New("repo err")

		userEmail = gofakeit.Email()
		user      = &models.User{
			ID:       gofakeit.Int64(),
			Email:    userEmail,
			Password: gofakeit.City(),
		}
		urlTitle = gofakeit.UUID()
		req      = &models.BuyTicketRequest{
			EventUrlTitle: urlTitle,
			FirstName:     gofakeit.FirstName(),
			LastName:      gofakeit.LastName(),
			PriceID:       gofakeit.Int64(),
			UserEmail:     userEmail,
		}
		event = &models.Event{
			ID:            gofakeit.Int64(),
			Title:         gofakeit.BeerName(),
			URLTitle:      urlTitle,
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
		ticket = &models.Ticket{
			UserID:    user.ID,
			User:      user,
			EventID:   event.ID,
			Event:     event,
			IsUsed:    false,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			PaymentID: "",
		}
	)
	closer.SetGlobalCloser(closer.New(wg))
	tests := []struct {
		name           string
		want           string
		err            error
		repositoryMock repositoryMockFunc
	}{{
		name: "success case",
		want: ticketID,
		err:  nil,
		repositoryMock: func(mc *minimock.Controller) repository.Repository {
			mock := mocks.NewRepositoryMock(mc)
			mock.UserMock.Expect(ctx, userEmail).Return(user, nil)
			mock.EventByURLTitleMock.Expect(ctx, urlTitle).Return(event, nil)
			// generates its own ticket id, idk how to handle this
			mock.InsertTicketMock.Expect(ctx, ticket).Return(ticketID, nil)
			return mock
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryMock := tt.repositoryMock(mc)

			service := apiService.NewApiService(wg, repositoryMock, authCfg, yooClient, mail)
			id, err := service.BuyTicket(ctx, req)

			require.Equal(t, tt.want, id)
			require.Equal(t, tt.err, err)
		})
	}
}

func TestTicket(t *testing.T) {
	t.Parallel()

	type repositoryMockFunc func(mc *minimock.Controller) repository.Repository

	var (
		wg  = &sync.WaitGroup{}
		ctx = context.Background()
		mc  = minimock.NewController(t)

		authCfg   = config.NewAuthConfig()
		yooCfg    = config.NewYookassaConfig()
		mailerCfg = config.NewMailerConfig()
		yooClient = yookassa.NewClient(yooCfg.ShopID(), yooCfg.ShopKey())
		mail, _   = smtp.NewSMTPMailer(mailerCfg, wg)

		repoErr = errors.New("repo err")

		ticketID = gofakeit.UUID()
		ticket   = &models.Ticket{
			ID:        ticketID,
			UserID:    gofakeit.Int64(),
			EventID:   gofakeit.Int64(),
			IsUsed:    false,
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
			PaymentID: "",
		}
	)
	closer.SetGlobalCloser(closer.New(wg))

	tests := []struct {
		name           string
		want           *models.Ticket
		err            error
		repositoryMock repositoryMockFunc
	}{
		{
			name: "success case",
			want: ticket,
			err:  nil,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.TicketMock.Expect(ctx, ticketID).Return(ticket, nil)
				return mock
			},
		},
		{
			name: "failure case",
			want: nil,
			err:  repoErr,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.TicketMock.Expect(ctx, ticketID).Return(nil, repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryMock := tt.repositoryMock(mc)

			service := apiService.NewApiService(wg, repositoryMock, authCfg, yooClient, mail)
			ticket, err := service.Ticket(ctx, ticketID)

			require.Equal(t, tt.want, ticket)
			require.Equal(t, tt.err, err)
		})
	}
}

func TestUserTickets(t *testing.T) {
	t.Parallel()

	type repositoryMockFunc func(mc *minimock.Controller) repository.Repository

	var (
		wg  = &sync.WaitGroup{}
		ctx = context.Background()
		mc  = minimock.NewController(t)

		authCfg   = config.NewAuthConfig()
		yooCfg    = config.NewYookassaConfig()
		mailerCfg = config.NewMailerConfig()
		yooClient = yookassa.NewClient(yooCfg.ShopID(), yooCfg.ShopKey())
		mail, _   = smtp.NewSMTPMailer(mailerCfg, wg)

		repoErr = errors.New("repo err")

		userID = gofakeit.Int64()

		ticketID = gofakeit.UUID()
		tickets  = []*models.Ticket{
			{
				ID:        ticketID,
				UserID:    userID,
				EventID:   gofakeit.Int64(),
				IsUsed:    false,
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				PaymentID: "",
			},
		}
	)
	closer.SetGlobalCloser(closer.New(wg))

	tests := []struct {
		name           string
		want           []*models.Ticket
		err            error
		repositoryMock repositoryMockFunc
	}{
		{
			name: "success case",
			want: tickets,
			err:  nil,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.UserTicketsMock.Expect(ctx, userID).Return(tickets, nil)
				return mock
			},
		},
		{
			name: "failure case",
			want: nil,
			err:  repoErr,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.UserTicketsMock.Expect(ctx, userID).Return(nil, repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryMock := tt.repositoryMock(mc)

			service := apiService.NewApiService(wg, repositoryMock, authCfg, yooClient, mail)
			ticket, err := service.UserTickets(ctx, userID)

			require.Equal(t, tt.want, ticket)
			require.Equal(t, tt.err, err)
		})
	}
}
