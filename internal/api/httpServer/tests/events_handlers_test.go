package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/wDRxxx/eventflow-backend/internal/api/httpServer"
	"github.com/wDRxxx/eventflow-backend/internal/config"
	"github.com/wDRxxx/eventflow-backend/internal/models"
	"github.com/wDRxxx/eventflow-backend/internal/oauth"
	"github.com/wDRxxx/eventflow-backend/internal/service"
	"github.com/wDRxxx/eventflow-backend/internal/service/mocks"
	"github.com/wDRxxx/eventflow-backend/internal/utils"
)

func TestEvents(t *testing.T) {
	t.Parallel()

	type apiServiceMockFunc func(mc *minimock.Controller) service.EventsService
	var (
		authCfg  = config.NewAuthConfig()
		httpCfg  = config.NewHttpConfig()
		oauthCfg = config.NewOAuthConfig()

		oauth = oauth.NewOAuth(oauthCfg)

		ctx  = context.Background()
		mc   = minimock.NewController(t)
		page = gofakeit.Int()

		method = http.MethodGet
		url    = "/api/events?page="

		serviceErr = errors.New("service error")

		events = []*models.Event{
			{
				Title:         gofakeit.BeerName(),
				URLTitle:      gofakeit.UUID(),
				Description:   gofakeit.ProductDescription(),
				BeginningTime: time.Now().Round(0),
				EndTime:       time.Now().Round(0),
				CreatorID:     gofakeit.Int64(),
				IsPublic:      false,
				IsFree:        true,
				PreviewImage:  gofakeit.UUID(),
				UTCOffset:     gofakeit.Int64(),
				Capacity:      gofakeit.Int64(),
				MinimalAge:    gofakeit.Int64(),
			},
		}
	)

	tests := []struct {
		name       string
		want       []*models.Event
		err        error
		statusCode int

		apiServiceMock apiServiceMockFunc
	}{
		{
			name:       "success case",
			want:       events,
			err:        nil,
			statusCode: http.StatusOK,
			apiServiceMock: func(mc *minimock.Controller) service.EventsService {
				mock := mocks.NewEventsServiceMock(mc)
				mock.EventsMock.Expect(minimock.AnyContext, page).Return(events, nil)
				return mock
			},
		},
		{
			name:       "negative page case",
			want:       nil,
			err:        nil,
			statusCode: http.StatusBadRequest,
			apiServiceMock: func(mc *minimock.Controller) service.EventsService {
				mock := mocks.NewEventsServiceMock(mc)
				return mock
			},
		},
		{
			name:       "failure case",
			want:       nil,
			err:        serviceErr,
			statusCode: http.StatusInternalServerError,
			apiServiceMock: func(mc *minimock.Controller) service.EventsService {
				mock := mocks.NewEventsServiceMock(mc)
				mock.EventsMock.Expect(minimock.AnyContext, page).Return(nil, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			apiServiceMock := tt.apiServiceMock(mc)

			api := httpServer.NewHTTPServer(
				authCfg,
				httpCfg,
				apiServiceMock,
				nil,
				nil,
				oauth,
			)
			server := httptest.NewServer(api.Handler())
			defer server.Close()

			var req *http.Request
			if tt.name == "negative page case" {
				req, _ = http.NewRequestWithContext(
					ctx,
					method,
					server.URL+url+fmt.Sprint(-page),
					nil,
				)
			} else {
				req, _ = http.NewRequestWithContext(
					ctx,
					method,
					server.URL+url+fmt.Sprint(page),
					nil,
				)
			}

			resp, _ := server.Client().Do(req)

			var respEvents []*models.Event
			_ = json.NewDecoder(resp.Body).Decode(&respEvents)

			require.Equal(t, tt.statusCode, resp.StatusCode)
			if respEvents != nil && len(respEvents) > 0 {
				for i, event := range respEvents {
					require.Equal(t, events[i], event)
				}
			}
		})
	}
}

func TestMyEvents(t *testing.T) {
	t.Parallel()

	type apiServiceMockFunc func(mc *minimock.Controller) service.EventsService

	var (
		authCfg  = config.NewAuthConfig()
		httpCfg  = config.NewHttpConfig()
		oauthCfg = config.NewOAuthConfig()

		oauth = oauth.NewOAuth(oauthCfg)

		ctx       = context.Background()
		mc        = minimock.NewController(t)
		creatorID = gofakeit.Int64()

		method = http.MethodGet
		url    = "/api/user/events"

		serviceErr = errors.New("service error")

		user = &models.UserClaims{
			RegisteredClaims: jwt.RegisteredClaims{Subject: fmt.Sprint(creatorID)},
			Email:            gofakeit.Email(),
		}

		events = []*models.Event{
			{
				Title:         gofakeit.BeerName(),
				URLTitle:      gofakeit.UUID(),
				Location:      gofakeit.ProgrammingLanguage(),
				Description:   gofakeit.ProductDescription(),
				BeginningTime: time.Now().Round(0),
				EndTime:       time.Now().Round(0),
				CreatorID:     creatorID,
				IsPublic:      false,
				IsFree:        true,
				PreviewImage:  gofakeit.UUID(),
				UTCOffset:     gofakeit.Int64(),
				Capacity:      gofakeit.Int64(),
				MinimalAge:    gofakeit.Int64(),
			},
		}
	)

	tests := []struct {
		name         string
		want         []*models.Event
		err          error
		isAuthorized bool
		statusCode   int

		apiServiceMock apiServiceMockFunc
	}{
		{
			name:         "unauthorized case",
			want:         nil,
			isAuthorized: false,
			statusCode:   http.StatusUnauthorized,
			apiServiceMock: func(mc *minimock.Controller) service.EventsService {
				mock := mocks.NewEventsServiceMock(mc)
				return mock
			},
		},
		{
			name:         "success case",
			want:         events,
			isAuthorized: true,
			err:          nil,
			statusCode:   http.StatusOK,
			apiServiceMock: func(mc *minimock.Controller) service.EventsService {
				mock := mocks.NewEventsServiceMock(mc)
				mock.UserEventsMock.Expect(minimock.AnyContext, creatorID).Return(events, nil)
				return mock
			},
		},
		{
			name:         "failure case",
			want:         nil,
			isAuthorized: true,
			err:          serviceErr,
			statusCode:   http.StatusInternalServerError,
			apiServiceMock: func(mc *minimock.Controller) service.EventsService {
				mock := mocks.NewEventsServiceMock(mc)
				mock.UserEventsMock.Expect(minimock.AnyContext, creatorID).Return(nil, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			apiServiceMock := tt.apiServiceMock(mc)

			api := httpServer.NewHTTPServer(
				authCfg,
				httpCfg,
				apiServiceMock,
				nil,
				nil,
				oauth,
			)

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

			var respEvents []*models.Event
			_ = json.NewDecoder(resp.Body).Decode(&respEvents)

			require.Equal(t, tt.statusCode, resp.StatusCode)
			if respEvents != nil && len(respEvents) > 0 {
				for i, event := range respEvents {
					require.Equal(t, events[i], event)
				}
			}
		})
	}
}

func TestEvent(t *testing.T) {
	t.Parallel()

	type apiServiceMockFunc func(mc *minimock.Controller) service.EventsService

	var (
		authCfg  = config.NewAuthConfig()
		httpCfg  = config.NewHttpConfig()
		oauthCfg = config.NewOAuthConfig()

		oauth = oauth.NewOAuth(oauthCfg)

		ctx      = context.Background()
		mc       = minimock.NewController(t)
		urlTitle = gofakeit.UUID()

		method = http.MethodGet
		url    = fmt.Sprintf("/api/events/%s", urlTitle)

		serviceErr = errors.New("service error")

		event = &models.Event{
			Title:         gofakeit.BeerName(),
			URLTitle:      urlTitle,
			Description:   gofakeit.ProductDescription(),
			BeginningTime: time.Now().Round(0),
			EndTime:       time.Now().Round(0),
			CreatorID:     gofakeit.Int64(),
			IsPublic:      false,
			IsFree:        true,
			PreviewImage:  gofakeit.UUID(),
			UTCOffset:     gofakeit.Int64(),
			Capacity:      gofakeit.Int64(),
			MinimalAge:    gofakeit.Int64(),
		}
	)

	tests := []struct {
		name       string
		want       *models.Event
		statusCode int

		apiServiceMock apiServiceMockFunc
	}{
		{
			name:       "success case",
			want:       event,
			statusCode: http.StatusOK,
			apiServiceMock: func(mc *minimock.Controller) service.EventsService {
				mock := mocks.NewEventsServiceMock(mc)
				mock.EventMock.Expect(minimock.AnyContext, urlTitle).Return(event, nil)
				return mock
			},
		},
		{
			name:       "not found case",
			want:       nil,
			statusCode: http.StatusNotFound,
			apiServiceMock: func(mc *minimock.Controller) service.EventsService {
				mock := mocks.NewEventsServiceMock(mc)
				mock.EventMock.Expect(minimock.AnyContext, urlTitle).Return(nil, pgx.ErrNoRows)
				return mock
			},
		},
		{
			name:       "failure case",
			want:       nil,
			statusCode: http.StatusInternalServerError,
			apiServiceMock: func(mc *minimock.Controller) service.EventsService {
				mock := mocks.NewEventsServiceMock(mc)
				mock.EventMock.Expect(minimock.AnyContext, urlTitle).Return(nil, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			apiServiceMock := tt.apiServiceMock(mc)

			api := httpServer.NewHTTPServer(
				authCfg,
				httpCfg,
				apiServiceMock,
				nil,
				nil,
				oauth,
			)

			server := httptest.NewServer(api.Handler())
			defer server.Close()
			req, _ := http.NewRequestWithContext(
				ctx,
				method,
				server.URL+url,
				nil,
			)

			resp, _ := server.Client().Do(req)

			require.Equal(t, tt.statusCode, resp.StatusCode)

			var respEvent *models.Event
			if tt.want != nil {
				_ = json.NewDecoder(resp.Body).Decode(&respEvent)
			}
			require.Equal(t, tt.want, respEvent)
		})
	}
}

func TestCreateEvent(t *testing.T) {
	t.Parallel()

	type apiServiceMockFunc func(mc *minimock.Controller) service.EventsService

	var (
		authCfg  = config.NewAuthConfig()
		httpCfg  = config.NewHttpConfig()
		oauthCfg = config.NewOAuthConfig()

		oauth = oauth.NewOAuth(oauthCfg)

		ctx       = context.Background()
		mc        = minimock.NewController(t)
		id        = gofakeit.Int64()
		creatorID = gofakeit.Int64()
		now       = time.Now().Round(0)

		method = http.MethodPost
		url    = "/api/events/"

		serviceErr = errors.New("service error")

		user = &models.UserClaims{
			RegisteredClaims: jwt.RegisteredClaims{Subject: fmt.Sprint(creatorID)},
			Email:            gofakeit.Email(),
		}

		event = &models.Event{
			Title:         gofakeit.BeerName(),
			URLTitle:      gofakeit.UUID(),
			Location:      gofakeit.ProgrammingLanguage(),
			Description:   gofakeit.ProductDescription(),
			BeginningTime: now,
			EndTime:       now,
			CreatorID:     creatorID,
			IsPublic:      false,
			IsFree:        true,
			PreviewImage:  gofakeit.UUID(),
			UTCOffset:     gofakeit.Int64(),
			Capacity:      gofakeit.Int64(),
			MinimalAge:    gofakeit.Int64(),
		}

		wrongEvent = &models.Event{
			URLTitle:      gofakeit.UUID(),
			Location:      gofakeit.ProgrammingLanguage(),
			Description:   gofakeit.ProductDescription(),
			BeginningTime: now,
			EndTime:       now,
			CreatorID:     creatorID,
			IsPublic:      false,
			IsFree:        true,
			PreviewImage:  gofakeit.UUID(),
			UTCOffset:     gofakeit.Int64(),
			Capacity:      gofakeit.Int64(),
			MinimalAge:    gofakeit.Int64(),
		}
	)

	tests := []struct {
		name         string
		event        *models.Event
		want         *models.DefaultResponse
		statusCode   int
		isAuthorized bool

		apiServiceMock apiServiceMockFunc
	}{
		{
			name:         "unauthorized case",
			want:         nil,
			isAuthorized: false,
			event:        event,
			statusCode:   http.StatusUnauthorized,
			apiServiceMock: func(mc *minimock.Controller) service.EventsService {
				mock := mocks.NewEventsServiceMock(mc)
				return mock
			},
		},
		{
			name: "success case",
			want: &models.DefaultResponse{
				Error:   false,
				Message: "created successfully",
			},
			isAuthorized: true,
			event:        event,
			statusCode:   http.StatusCreated,
			apiServiceMock: func(mc *minimock.Controller) service.EventsService {
				mock := mocks.NewEventsServiceMock(mc)
				mock.CreateEventMock.Expect(minimock.AnyContext, event).Return(id, nil)
				return mock
			},
		},
		{
			name:         "failure case",
			want:         nil,
			isAuthorized: true,
			event:        event,
			statusCode:   http.StatusInternalServerError,
			apiServiceMock: func(mc *minimock.Controller) service.EventsService {
				mock := mocks.NewEventsServiceMock(mc)
				mock.CreateEventMock.Expect(minimock.AnyContext, event).Return(0, serviceErr)
				return mock
			},
		},
		{
			name:         "wrong input case",
			want:         nil,
			isAuthorized: true,
			event:        wrongEvent,
			statusCode:   http.StatusUnprocessableEntity,
			apiServiceMock: func(mc *minimock.Controller) service.EventsService {
				mock := mocks.NewEventsServiceMock(mc)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			apiServiceMock := tt.apiServiceMock(mc)

			api := httpServer.NewHTTPServer(
				authCfg,
				httpCfg,
				apiServiceMock,
				nil,
				nil,
				oauth,
			)

			server := httptest.NewServer(api.Handler())
			defer server.Close()

			data, _ := json.Marshal(tt.event)
			form := map[string]string{
				"image": "",
				"event": string(data),
			}

			contentType, body, _ := createForm(form)
			req, _ := http.NewRequestWithContext(
				ctx,
				method,
				server.URL+url,
				body,
			)

			req.Header.Set("Content-Type", contentType)
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

func TestUpdateEvent(t *testing.T) {
	t.Parallel()

	type apiServiceMockFunc func(mc *minimock.Controller) service.EventsService

	var (
		authCfg  = config.NewAuthConfig()
		httpCfg  = config.NewHttpConfig()
		oauthCfg = config.NewOAuthConfig()

		oauth = oauth.NewOAuth(oauthCfg)

		ctx       = context.Background()
		mc        = minimock.NewController(t)
		creatorID = gofakeit.Int64()
		now       = time.Now().Round(0)
		urlTitle  = gofakeit.UUID()

		method = http.MethodPut
		url    = fmt.Sprintf("/api/events/%s", urlTitle)

		serviceErr = errors.New("service error")

		user = &models.UserClaims{
			RegisteredClaims: jwt.RegisteredClaims{Subject: fmt.Sprint(creatorID)},
			Email:            gofakeit.Email(),
		}

		event = &models.Event{
			Title:         gofakeit.BeerName(),
			URLTitle:      urlTitle,
			Location:      gofakeit.ProgrammingLanguage(),
			Description:   gofakeit.ProductDescription(),
			BeginningTime: now,
			EndTime:       now,
			CreatorID:     creatorID,
			IsPublic:      false,
			IsFree:        true,
			PreviewImage:  gofakeit.UUID(),
			UTCOffset:     gofakeit.Int64(),
			Capacity:      gofakeit.Int64(),
			MinimalAge:    gofakeit.Int64(),
		}
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
			isAuthorized: false,
			statusCode:   http.StatusUnauthorized,
			apiServiceMock: func(mc *minimock.Controller) service.EventsService {
				mock := mocks.NewEventsServiceMock(mc)
				return mock
			},
		},
		{
			name: "success case",
			want: &models.DefaultResponse{
				Error:   false,
				Message: "updated successfully",
			},
			isAuthorized: true,
			statusCode:   http.StatusOK,
			apiServiceMock: func(mc *minimock.Controller) service.EventsService {
				mock := mocks.NewEventsServiceMock(mc)
				mock.UpdateEventMock.Expect(minimock.AnyContext, creatorID, event).Return(nil)
				return mock
			},
		},
		{
			name:         "failure case",
			want:         nil,
			isAuthorized: true,
			statusCode:   http.StatusInternalServerError,
			apiServiceMock: func(mc *minimock.Controller) service.EventsService {
				mock := mocks.NewEventsServiceMock(mc)
				mock.UpdateEventMock.Expect(minimock.AnyContext, creatorID, event).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			apiServiceMock := tt.apiServiceMock(mc)

			api := httpServer.NewHTTPServer(
				authCfg,
				httpCfg,
				apiServiceMock,
				nil,
				nil,
				oauth,
			)

			server := httptest.NewServer(api.Handler())
			defer server.Close()

			data, _ := json.Marshal(event)
			form := map[string]string{
				"image": "",
				"event": string(data),
			}

			contentType, body, _ := createForm(form)
			req, _ := http.NewRequestWithContext(
				ctx,
				method,
				server.URL+url,
				body,
			)

			req.Header.Set("Content-Type", contentType)
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

func TestDeleteEvent(t *testing.T) {
	t.Parallel()

	type apiServiceMockFunc func(mc *minimock.Controller) service.EventsService

	var (
		authCfg  = config.NewAuthConfig()
		httpCfg  = config.NewHttpConfig()
		oauthCfg = config.NewOAuthConfig()

		oauth = oauth.NewOAuth(oauthCfg)

		ctx       = context.Background()
		mc        = minimock.NewController(t)
		creatorID = gofakeit.Int64()
		now       = time.Now().Round(0)
		urlTitle  = gofakeit.UUID()

		method = http.MethodDelete
		url    = fmt.Sprintf("/api/events/%s", urlTitle)

		serviceErr = errors.New("service error")

		user = &models.UserClaims{
			RegisteredClaims: jwt.RegisteredClaims{Subject: fmt.Sprint(creatorID)},
			Email:            gofakeit.Email(),
		}

		event = &models.Event{
			Title:         gofakeit.BeerName(),
			URLTitle:      urlTitle,
			Location:      gofakeit.ProgrammingLanguage(),
			Description:   gofakeit.ProductDescription(),
			BeginningTime: now,
			EndTime:       now,
			CreatorID:     creatorID,
			IsPublic:      false,
			IsFree:        true,
			PreviewImage:  gofakeit.UUID(),
			UTCOffset:     gofakeit.Int64(),
			Capacity:      gofakeit.Int64(),
			MinimalAge:    gofakeit.Int64(),
		}
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
			isAuthorized: false,
			statusCode:   http.StatusUnauthorized,
			apiServiceMock: func(mc *minimock.Controller) service.EventsService {
				mock := mocks.NewEventsServiceMock(mc)
				return mock
			},
		},
		{
			name: "success case",
			want: &models.DefaultResponse{
				Error:   false,
				Message: "Event was deleted successfully",
			},
			isAuthorized: true,
			statusCode:   http.StatusOK,
			apiServiceMock: func(mc *minimock.Controller) service.EventsService {
				mock := mocks.NewEventsServiceMock(mc)
				mock.DeleteEventMock.Expect(minimock.AnyContext, creatorID, urlTitle).Return(nil)
				return mock
			},
		},
		{
			name:         "failure case",
			want:         nil,
			isAuthorized: true,
			statusCode:   http.StatusInternalServerError,
			apiServiceMock: func(mc *minimock.Controller) service.EventsService {
				mock := mocks.NewEventsServiceMock(mc)
				mock.DeleteEventMock.Expect(minimock.AnyContext, creatorID, urlTitle).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			apiServiceMock := tt.apiServiceMock(mc)

			api := httpServer.NewHTTPServer(
				authCfg,
				httpCfg,
				apiServiceMock,
				nil,
				nil,
				oauth,
			)

			server := httptest.NewServer(api.Handler())
			defer server.Close()

			data, _ := json.Marshal(event)
			form := map[string]string{
				"image": "",
				"event": string(data),
			}

			contentType, body, _ := createForm(form)
			req, _ := http.NewRequestWithContext(
				ctx,
				method,
				server.URL+url,
				body,
			)

			req.Header.Set("Content-Type", contentType)
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
