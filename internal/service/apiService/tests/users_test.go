package tests

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/wDRxxx/eventflow-backend/internal/closer"
	"github.com/wDRxxx/eventflow-backend/internal/config"
	"github.com/wDRxxx/eventflow-backend/internal/mailer/smtp"
	"github.com/wDRxxx/eventflow-backend/internal/models"
	"github.com/wDRxxx/eventflow-backend/internal/repository"
	"github.com/wDRxxx/eventflow-backend/internal/repository/mocks"
	"github.com/wDRxxx/eventflow-backend/internal/service"
	"github.com/wDRxxx/eventflow-backend/internal/service/apiService"
	"github.com/wDRxxx/eventflow-backend/internal/utils"
)

func TestRegisterUser(t *testing.T) {
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

		userID    = gofakeit.Int64()
		userEmail = gofakeit.Email()
		user      = &models.User{
			ID:    userID,
			Email: userEmail,
		}
	)
	closer.SetGlobalCloser(closer.New(wg))

	tests := []struct {
		name           string
		err            error
		repositoryMock repositoryMockFunc
	}{
		{
			name: "success case",
			err:  nil,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.UserMock.Expect(ctx, userEmail).Return(nil, repoErr)
				mock.InsertUserMock.Expect(ctx, user).Return(userID, nil)
				return mock
			},
		},
		{
			name: "user exists case",
			err:  service.ErrUserAlreadyExists,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.UserMock.Expect(ctx, userEmail).Return(user, nil)
				return mock
			},
		},
		{
			name: "failure case",
			err:  repoErr,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.UserMock.Expect(ctx, userEmail).Return(nil, repoErr)
				mock.InsertUserMock.Expect(ctx, user).Return(0, repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryMock := tt.repositoryMock(mc)

			service := apiService.NewApiService(wg, repositoryMock, authCfg, mail)
			err := service.RegisterUser(ctx, user)

			require.Equal(t, tt.err, err)
		})
	}
}

func TestLogin(t *testing.T) {
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

		userID    = gofakeit.Int64()
		userEmail = gofakeit.Email()

		pass        = gofakeit.Password(true, true, true, false, false, 12)
		hashPass, _ = bcrypt.GenerateFromPassword([]byte(pass), 12)

		user = &models.User{
			ID:       userID,
			Email:    userEmail,
			Password: pass,
		}
		dbUser = &models.User{
			ID:       userID,
			Email:    userEmail,
			Password: string(hashPass),
		}
	)
	closer.SetGlobalCloser(closer.New(wg))

	tests := []struct {
		name           string
		err            error
		repositoryMock repositoryMockFunc
	}{
		{
			name: "success case",
			err:  nil,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.UserMock.Expect(ctx, userEmail).Return(dbUser, nil)
				return mock
			},
		},
		{
			name: "user doesn't exists case",
			err:  service.ErrWrongCredentials,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.UserMock.Expect(ctx, userEmail).Return(nil, repoErr)
				return mock
			},
		},
		{
			name: "wrong password case",
			err:  service.ErrWrongCredentials,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				user := user
				user.Password = "z"
				mock.UserMock.Expect(ctx, userEmail).Return(user, repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryMock := tt.repositoryMock(mc)

			service := apiService.NewApiService(wg, repositoryMock, authCfg, mail)
			_, err := service.Login(ctx, user)

			require.Equal(t, tt.err, err)
		})
	}
}

func TestAccessToken(t *testing.T) {
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

		userID    = gofakeit.Int64()
		userEmail = gofakeit.Email()
		user      = &models.User{
			ID:    userID,
			Email: userEmail,
		}
	)
	closer.SetGlobalCloser(closer.New(wg))

	refreshToken, _ := utils.GenerateToken(&models.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: fmt.Sprint(user.ID),
		},
		Email: userEmail,
	}, authCfg.RefreshTokenSecret(), authCfg.RefreshTokenTTL())

	tests := []struct {
		name           string
		token          string
		err            error
		repositoryMock repositoryMockFunc
	}{
		{
			name:  "success case",
			token: refreshToken,
			err:   nil,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				return mock
			},
		},
		{
			name:  "failure case",
			token: "",
			err:   repoErr,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryMock := tt.repositoryMock(mc)

			service := apiService.NewApiService(wg, repositoryMock, authCfg, mail)
			_, err := service.AccessToken(ctx, tt.token)

			if tt.err != nil {
				require.Error(t, err)
			} else {
				require.Equal(t, tt.err, err)
			}
		})
	}
}

func TestUser(t *testing.T) {
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

		userID    = gofakeit.Int64()
		userEmail = gofakeit.Email()

		pass        = gofakeit.Password(true, true, true, false, false, 12)
		hashPass, _ = bcrypt.GenerateFromPassword([]byte(pass), 12)

		dbUser = &models.User{
			ID:       userID,
			Email:    userEmail,
			Password: string(hashPass),
		}
	)
	closer.SetGlobalCloser(closer.New(wg))

	tests := []struct {
		name           string
		want           *models.User
		err            error
		repositoryMock repositoryMockFunc
	}{
		{
			name: "success case",
			want: dbUser,
			err:  nil,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.UserMock.Expect(ctx, userEmail).Return(dbUser, nil)
				return mock
			},
		},
		{
			name: "failure case",
			want: nil,
			err:  repoErr,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.UserMock.Expect(ctx, userEmail).Return(nil, repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryMock := tt.repositoryMock(mc)

			service := apiService.NewApiService(wg, repositoryMock, authCfg, mail)
			u, err := service.User(ctx, userEmail)

			require.Equal(t, tt.want, u)
			require.Equal(t, tt.err, err)
		})
	}
}

func TestUpdateUser(t *testing.T) {
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

		userID    = gofakeit.Int64()
		userEmail = gofakeit.Email()
		username  = gofakeit.Username()

		settings = models.YookassaSettings{
			UserID:  userID,
			ShopID:  gofakeit.BeerName(),
			ShopKey: gofakeit.BeerYeast(),
		}

		user = &models.User{
			ID:               userID,
			Email:            userEmail,
			TGUsername:       username,
			YookassaSettings: settings,
		}
	)
	closer.SetGlobalCloser(closer.New(wg))

	tests := []struct {
		name           string
		err            error
		repositoryMock repositoryMockFunc
	}{
		{
			name: "success case",
			err:  nil,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.UpdateUserTGUsernameMock.Expect(ctx, userID, username).Return(nil)
				mock.UpdateYookassaSettingsMock.Expect(ctx, &settings).Return(nil)
				return mock
			},
		},
		{
			name: "failure case 1",
			err:  repoErr,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.UpdateUserTGUsernameMock.Expect(ctx, userID, username).Return(repoErr)
				return mock
			},
		},
		{
			name: "failure case 2",
			err:  repoErr,
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.UpdateUserTGUsernameMock.Expect(ctx, userID, username).Return(nil)
				mock.UpdateYookassaSettingsMock.Expect(ctx, &settings).Return(repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryMock := tt.repositoryMock(mc)

			service := apiService.NewApiService(wg, repositoryMock, authCfg, mail)
			err := service.UpdateUser(ctx, user)

			require.Equal(t, tt.err, err)
		})
	}
}
