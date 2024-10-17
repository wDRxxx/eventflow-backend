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

func TestRegister(t *testing.T) {
	t.Parallel()

	type apiServiceMockFunc func(mc *minimock.Controller) service.UsersService

	var (
		authCfg = config.NewAuthConfig()
		httpCfg = config.NewHttpConfig()

		ctx = context.Background()
		mc  = minimock.NewController(t)

		method = http.MethodPost
		url    = "/api/auth/register"

		serviceErr = errors.New("service error")

		user = &models.User{
			Email:    gofakeit.Email(),
			Password: gofakeit.Password(true, true, true, false, false, 12),
		}
		wrongUser = &models.User{
			Email:    gofakeit.Email(),
			Password: gofakeit.Password(true, true, true, false, false, 12),
		}
	)

	tests := []struct {
		name       string
		want       *models.DefaultResponse
		statusCode int
		user       *models.User

		apiServiceMock apiServiceMockFunc
	}{
		{
			name: "success case",
			want: &models.DefaultResponse{
				Error:   false,
				Message: "You was successfully registered",
			},
			user:       user,
			statusCode: http.StatusCreated,
			apiServiceMock: func(mc *minimock.Controller) service.UsersService {
				mock := mocks.NewUsersServiceMock(mc)
				mock.RegisterUserMock.Expect(minimock.AnyContext, user).Return(nil)
				return mock
			},
		},
		{
			name:       "failure case",
			want:       nil,
			user:       user,
			statusCode: http.StatusInternalServerError,
			apiServiceMock: func(mc *minimock.Controller) service.UsersService {
				mock := mocks.NewUsersServiceMock(mc)
				mock.RegisterUserMock.Expect(minimock.AnyContext, user).Return(serviceErr)
				return mock
			},
		},
		{
			name:       "wrong input case",
			want:       nil,
			user:       wrongUser,
			statusCode: http.StatusInternalServerError,
			apiServiceMock: func(mc *minimock.Controller) service.UsersService {
				mock := mocks.NewUsersServiceMock(mc)
				mock.RegisterUserMock.Expect(minimock.AnyContext, wrongUser).Return(serviceErr)
				return mock
			},
		},
		{
			name:       "wrong input case",
			want:       nil,
			user:       wrongUser,
			statusCode: http.StatusConflict,
			apiServiceMock: func(mc *minimock.Controller) service.UsersService {
				mock := mocks.NewUsersServiceMock(mc)
				mock.RegisterUserMock.Expect(minimock.AnyContext, wrongUser).Return(service.ErrUserAlreadyExists)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServiceMock := tt.apiServiceMock(mc)

			api := httpServer.NewHTTPServer(
				authCfg,
				httpCfg,
				nil,
				nil,
				apiServiceMock,
			)

			server := httptest.NewServer(api.Handler())
			defer server.Close()

			data, _ := json.Marshal(tt.user)

			buf := bytes.NewBuffer(data)

			req, _ := http.NewRequestWithContext(
				ctx,
				method,
				server.URL+url,
				buf,
			)

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

func TestLogin(t *testing.T) {
	t.Parallel()

	type apiServiceMockFunc func(mc *minimock.Controller) service.UsersService

	var (
		authCfg = config.NewAuthConfig()
		httpCfg = config.NewHttpConfig()

		ctx = context.Background()
		mc  = minimock.NewController(t)

		method = http.MethodPost
		url    = "/api/auth/login"

		serviceErr = errors.New("service error")

		user = &models.User{
			Email:    gofakeit.Email(),
			Password: gofakeit.Password(true, true, true, false, false, 12),
		}
		wrongUser = &models.User{
			Email:    gofakeit.Email(),
			Password: gofakeit.Password(true, true, true, false, false, 12),
		}

		refreshToken = "refresh"
		accessToken  = "access"
	)

	tests := []struct {
		name       string
		want       *models.TokenPair
		statusCode int
		user       *models.User

		apiServiceMock apiServiceMockFunc
	}{
		{
			name: "success case",
			want: &models.TokenPair{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			},
			user:       user,
			statusCode: http.StatusOK,
			apiServiceMock: func(mc *minimock.Controller) service.UsersService {
				mock := mocks.NewUsersServiceMock(mc)
				mock.LoginMock.Expect(minimock.AnyContext, user).Return(refreshToken, nil)
				mock.AccessTokenMock.Expect(minimock.AnyContext, refreshToken).Return(accessToken, nil)
				return mock
			},
		},
		{
			name:       "failure case",
			want:       nil,
			user:       user,
			statusCode: http.StatusInternalServerError,
			apiServiceMock: func(mc *minimock.Controller) service.UsersService {
				mock := mocks.NewUsersServiceMock(mc)
				mock.LoginMock.Expect(minimock.AnyContext, user).Return("", serviceErr)
				return mock
			},
		},
		{
			name:       "wrong input case",
			want:       nil,
			user:       wrongUser,
			statusCode: http.StatusUnauthorized,
			apiServiceMock: func(mc *minimock.Controller) service.UsersService {
				mock := mocks.NewUsersServiceMock(mc)
				mock.LoginMock.Expect(minimock.AnyContext, wrongUser).Return("", service.ErrWrongCredentials)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServiceMock := tt.apiServiceMock(mc)

			api := httpServer.NewHTTPServer(
				authCfg,
				httpCfg,
				nil,
				nil,
				apiServiceMock,
			)

			server := httptest.NewServer(api.Handler())
			defer server.Close()

			data, _ := json.Marshal(tt.user)

			buf := bytes.NewBuffer(data)

			req, _ := http.NewRequestWithContext(
				ctx,
				method,
				server.URL+url,
				buf,
			)

			resp, _ := server.Client().Do(req)
			require.Equal(t, tt.statusCode, resp.StatusCode)
			var res *models.TokenPair
			if tt.want != nil {
				_ = json.NewDecoder(resp.Body).Decode(&res)
			}
			require.Equal(t, tt.want, res)
		})
	}
}

func TestRefresh(t *testing.T) {
	t.Parallel()

	type apiServiceMockFunc func(mc *minimock.Controller) service.UsersService

	var (
		authCfg = config.NewAuthConfig()
		httpCfg = config.NewHttpConfig()

		ctx = context.Background()
		mc  = minimock.NewController(t)

		method = http.MethodPost
		url    = "/api/auth/refresh"

		serviceErr = errors.New("service error")

		refreshToken = "refresh"
		accessToken  = "access"
	)

	tests := []struct {
		name       string
		want       *models.DefaultResponse
		statusCode int
		hasCookie  bool

		apiServiceMock apiServiceMockFunc
	}{
		{
			name: "success case",
			want: &models.DefaultResponse{
				Error:   false,
				Message: accessToken,
			},
			hasCookie:  true,
			statusCode: http.StatusOK,
			apiServiceMock: func(mc *minimock.Controller) service.UsersService {
				mock := mocks.NewUsersServiceMock(mc)
				mock.AccessTokenMock.Expect(minimock.AnyContext, refreshToken).Return(accessToken, nil)
				return mock
			},
		},
		{
			name:       "failure case",
			want:       nil,
			hasCookie:  true,
			statusCode: http.StatusInternalServerError,
			apiServiceMock: func(mc *minimock.Controller) service.UsersService {
				mock := mocks.NewUsersServiceMock(mc)
				mock.AccessTokenMock.Expect(minimock.AnyContext, refreshToken).Return("", serviceErr)
				return mock
			},
		},
		{
			name:       "no refresh token cookie case",
			want:       nil,
			hasCookie:  false,
			statusCode: http.StatusUnauthorized,
			apiServiceMock: func(mc *minimock.Controller) service.UsersService {
				mock := mocks.NewUsersServiceMock(mc)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServiceMock := tt.apiServiceMock(mc)

			api := httpServer.NewHTTPServer(
				authCfg,
				httpCfg,
				nil,
				nil,
				apiServiceMock,
			)

			server := httptest.NewServer(api.Handler())
			defer server.Close()

			req, _ := http.NewRequestWithContext(
				ctx,
				method,
				server.URL+url,
				nil,
			)
			if tt.hasCookie {
				req.AddCookie(&http.Cookie{
					Name:  "refresh_token",
					Value: refreshToken,
				})
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

func TestLogout(t *testing.T) {
	t.Parallel()

	type apiServiceMockFunc func(mc *minimock.Controller) service.UsersService

	var (
		authCfg = config.NewAuthConfig()
		httpCfg = config.NewHttpConfig()

		ctx = context.Background()
		mc  = minimock.NewController(t)

		method = http.MethodPost
		url    = "/api/auth/logout"
	)

	tests := []struct {
		name       string
		want       string
		statusCode int

		apiServiceMock apiServiceMockFunc
	}{
		{
			name:       "success case",
			want:       "",
			statusCode: http.StatusAccepted,
			apiServiceMock: func(mc *minimock.Controller) service.UsersService {
				mock := mocks.NewUsersServiceMock(mc)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServiceMock := tt.apiServiceMock(mc)

			api := httpServer.NewHTTPServer(
				authCfg,
				httpCfg,
				nil,
				nil,
				apiServiceMock,
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
			require.Equal(t, "", resp.Cookies()[0].Value)
		})
	}
}

func TestProfile(t *testing.T) {
	t.Parallel()

	type apiServiceMockFunc func(mc *minimock.Controller) service.UsersService

	var (
		authCfg = config.NewAuthConfig()
		httpCfg = config.NewHttpConfig()

		ctx = context.Background()
		mc  = minimock.NewController(t)

		method = http.MethodGet
		url    = "/api/user/profile"

		serviceErr = errors.New("service error")

		userID    = gofakeit.Int64()
		userEmail = gofakeit.Email()

		userClaims = &models.UserClaims{
			RegisteredClaims: jwt.RegisteredClaims{Subject: fmt.Sprint(userID)},
			Email:            userEmail,
		}
		user = &models.User{
			TGUsername: gofakeit.Username(),
			Email:      userEmail,
		}
	)

	tests := []struct {
		name         string
		want         *models.User
		isAuthorized bool
		statusCode   int

		apiServiceMock apiServiceMockFunc
	}{
		{
			name:         "success case",
			want:         user,
			statusCode:   http.StatusOK,
			isAuthorized: true,
			apiServiceMock: func(mc *minimock.Controller) service.UsersService {
				mock := mocks.NewUsersServiceMock(mc)
				mock.UserMock.Expect(minimock.AnyContext, userEmail).Return(user, nil)
				return mock
			},
		},
		{
			name:         "failure case",
			want:         nil,
			statusCode:   http.StatusInternalServerError,
			isAuthorized: true,
			apiServiceMock: func(mc *minimock.Controller) service.UsersService {
				mock := mocks.NewUsersServiceMock(mc)
				mock.UserMock.Expect(minimock.AnyContext, userEmail).Return(nil, serviceErr)
				return mock
			},
		},
		{
			name:         "unauthorized case",
			want:         nil,
			statusCode:   http.StatusUnauthorized,
			isAuthorized: false,
			apiServiceMock: func(mc *minimock.Controller) service.UsersService {
				mock := mocks.NewUsersServiceMock(mc)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServiceMock := tt.apiServiceMock(mc)

			api := httpServer.NewHTTPServer(
				authCfg,
				httpCfg,
				nil,
				nil,
				apiServiceMock,
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
				token, _ := utils.GenerateToken(userClaims, authCfg.AccessTokenSecret(), authCfg.AccessTokenTTL())
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			}

			resp, _ := server.Client().Do(req)
			require.Equal(t, tt.statusCode, resp.StatusCode)
			var res *models.User
			if tt.want != nil {
				_ = json.NewDecoder(resp.Body).Decode(&res)
			}
			require.Equal(t, tt.want, res)
		})
	}
}

func TestUpdateProfile(t *testing.T) {
	t.Parallel()

	type apiServiceMockFunc func(mc *minimock.Controller) service.UsersService

	var (
		authCfg = config.NewAuthConfig()
		httpCfg = config.NewHttpConfig()

		ctx = context.Background()
		mc  = minimock.NewController(t)

		method = http.MethodPut
		url    = "/api/user/profile"

		serviceErr = errors.New("service error")

		userID    = gofakeit.Int64()
		userEmail = gofakeit.Email()

		userClaims = &models.UserClaims{
			RegisteredClaims: jwt.RegisteredClaims{Subject: fmt.Sprint(userID)},
			Email:            userEmail,
		}
		user = &models.User{
			ID:         userID,
			TGUsername: gofakeit.Username(),
			Email:      userEmail,
		}
	)

	tests := []struct {
		name         string
		want         *models.DefaultResponse
		isAuthorized bool
		statusCode   int

		apiServiceMock apiServiceMockFunc
	}{
		{
			name: "success case",
			want: &models.DefaultResponse{
				Error:   false,
				Message: "Successfully updated user",
			},
			statusCode:   http.StatusOK,
			isAuthorized: true,
			apiServiceMock: func(mc *minimock.Controller) service.UsersService {
				mock := mocks.NewUsersServiceMock(mc)
				mock.UpdateUserMock.Expect(minimock.AnyContext, user).Return(nil)
				return mock
			},
		},
		{
			name:         "failure case",
			want:         nil,
			statusCode:   http.StatusInternalServerError,
			isAuthorized: true,
			apiServiceMock: func(mc *minimock.Controller) service.UsersService {
				mock := mocks.NewUsersServiceMock(mc)
				mock.UpdateUserMock.Expect(minimock.AnyContext, user).Return(serviceErr)
				return mock
			},
		},
		{
			name:         "unauthorized case",
			want:         nil,
			statusCode:   http.StatusUnauthorized,
			isAuthorized: false,
			apiServiceMock: func(mc *minimock.Controller) service.UsersService {
				mock := mocks.NewUsersServiceMock(mc)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServiceMock := tt.apiServiceMock(mc)

			api := httpServer.NewHTTPServer(
				authCfg,
				httpCfg,
				nil,
				nil,
				apiServiceMock,
			)

			server := httptest.NewServer(api.Handler())
			defer server.Close()

			data, _ := json.Marshal(user)
			buf := bytes.NewBuffer(data)

			req, _ := http.NewRequestWithContext(
				ctx,
				method,
				server.URL+url,
				buf,
			)
			if tt.isAuthorized {
				token, _ := utils.GenerateToken(userClaims, authCfg.AccessTokenSecret(), authCfg.AccessTokenTTL())
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
