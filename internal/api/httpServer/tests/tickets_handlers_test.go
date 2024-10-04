package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/wDRxxx/eventflow-backend/internal/api/httpServer"
	"github.com/wDRxxx/eventflow-backend/internal/config"
	"github.com/wDRxxx/eventflow-backend/internal/models"
	"github.com/wDRxxx/eventflow-backend/internal/service"
	"github.com/wDRxxx/eventflow-backend/internal/service/mocks"
	"github.com/wDRxxx/eventflow-backend/internal/utils"
)

func TestUserTickets(t *testing.T) {
	t.Parallel()

	type apiServiceMockFunc func(mc *minimock.Controller) service.ApiService

	var (
		authCfg = config.NewAuthConfig()
		httpCfg = config.NewHttpConfig()

		ctx = context.Background()
		mc  = minimock.NewController(t)

		method = http.MethodGet
		url    = "/api/user/tickets"

		serviceErr = errors.New("service error")

		userID = gofakeit.Int64()
		user   = &models.UserClaims{
			RegisteredClaims: jwt.RegisteredClaims{Subject: fmt.Sprint(userID)},
			Email:            gofakeit.Email(),
		}
		tickets = []*models.Ticket{
			{
				ID:        gofakeit.UUID(),
				IsUsed:    false,
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
			},
		}
	)

	tests := []struct {
		name         string
		want         []*models.Ticket
		statusCode   int
		isAuthorized bool

		apiServiceMock apiServiceMockFunc
	}{
		{
			name:         "unauthorized case",
			want:         nil,
			statusCode:   http.StatusUnauthorized,
			isAuthorized: false,
			apiServiceMock: func(mc *minimock.Controller) service.ApiService {
				mock := mocks.NewApiServiceMock(mc)
				return mock
			},
		},
		{
			name:         "success case",
			want:         tickets,
			statusCode:   http.StatusOK,
			isAuthorized: true,
			apiServiceMock: func(mc *minimock.Controller) service.ApiService {
				mock := mocks.NewApiServiceMock(mc)
				mock.UserTicketsMock.Expect(minimock.AnyContext, userID).Return(tickets, nil)
				return mock
			},
		},
		{
			name:         "failure case",
			want:         nil,
			statusCode:   http.StatusInternalServerError,
			isAuthorized: true,
			apiServiceMock: func(mc *minimock.Controller) service.ApiService {
				mock := mocks.NewApiServiceMock(mc)
				mock.UserTicketsMock.Expect(minimock.AnyContext, userID).Return(nil, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			apiServiceMock := tt.apiServiceMock(mc)

			api := httpServer.NewHTTPServer(
				apiServiceMock,
				authCfg,
				httpCfg)

			server := httptest.NewServer(api.Handler())
			defer server.Close()

			req, _ := http.NewRequestWithContext(
				ctx,
				method,
				server.URL+url,
				nil,
			)

			if tt.isAuthorized {
				token, _ := utils.GenerateToken(user, authCfg.AccessTokenSecret(), authCfg.AccessTokenTTL())
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			}

			resp, _ := server.Client().Do(req)
			require.Equal(t, tt.statusCode, resp.StatusCode)
			var res []*models.Ticket
			if tt.want != nil {
				_ = json.NewDecoder(resp.Body).Decode(&res)
			}
			require.Equal(t, tt.want, res)
		})
	}
}

func TestBuyTicket(t *testing.T) {
	t.Parallel()

	type apiServiceMockFunc func(mc *minimock.Controller) service.ApiService

	var (
		authCfg = config.NewAuthConfig()
		httpCfg = config.NewHttpConfig()

		ctx = context.Background()
		mc  = minimock.NewController(t)

		method = http.MethodPost
		url    = "/api/tickets"

		serviceErr = errors.New("service error")

		userID    = gofakeit.Int64()
		userEmail = gofakeit.Email()

		user = &models.UserClaims{
			RegisteredClaims: jwt.RegisteredClaims{Subject: fmt.Sprint(userID)},
			Email:            userEmail,
		}
		request = &models.BuyTicketRequest{
			EventUrlTitle: gofakeit.UUID(),
			FirstName:     gofakeit.FirstName(),
			LastName:      gofakeit.LastName(),
			PriceID:       gofakeit.Int64(),
			UserEmail:     userEmail,
		}
		returnURL = gofakeit.URL()
	)

	tests := []struct {
		name         string
		want         *models.DefaultResponse
		statusCode   int
		isAuthorized bool

		apiServiceMock apiServiceMockFunc
	}{
		{
			name:         "unauthorized case",
			want:         nil,
			statusCode:   http.StatusUnauthorized,
			isAuthorized: false,
			apiServiceMock: func(mc *minimock.Controller) service.ApiService {
				mock := mocks.NewApiServiceMock(mc)
				return mock
			},
		},
		{
			name: "success case",
			want: &models.DefaultResponse{
				Error:   false,
				Message: returnURL,
			},
			statusCode:   http.StatusOK,
			isAuthorized: true,
			apiServiceMock: func(mc *minimock.Controller) service.ApiService {
				mock := mocks.NewApiServiceMock(mc)
				mock.BuyTicketMock.Expect(minimock.AnyContext, request).Return(returnURL, nil)
				return mock
			},
		},
		{
			name:         "failure case",
			want:         nil,
			statusCode:   http.StatusInternalServerError,
			isAuthorized: true,
			apiServiceMock: func(mc *minimock.Controller) service.ApiService {
				mock := mocks.NewApiServiceMock(mc)
				mock.BuyTicketMock.Expect(minimock.AnyContext, request).Return("", serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			apiServiceMock := tt.apiServiceMock(mc)

			api := httpServer.NewHTTPServer(
				apiServiceMock,
				authCfg,
				httpCfg)

			server := httptest.NewServer(api.Handler())
			defer server.Close()

			data, _ := json.Marshal(request)

			buf := bytes.NewBuffer(data)

			req, _ := http.NewRequestWithContext(
				ctx,
				method,
				server.URL+url,
				buf,
			)

			if tt.isAuthorized {
				token, _ := utils.GenerateToken(user, authCfg.AccessTokenSecret(), authCfg.AccessTokenTTL())
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			}

			resp, _ := server.Client().Do(req)
			require.Equal(t, tt.statusCode, resp.StatusCode)
			var res *models.DefaultResponse
			if tt.want != nil {
				_ = json.NewDecoder(resp.Body).Decode(&res)
			}
			require.Equal(t, tt.want, res)
		})
	}
}
