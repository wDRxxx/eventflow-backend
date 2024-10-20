package usersService

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/wDRxxx/eventflow-backend/internal/config"
	"github.com/wDRxxx/eventflow-backend/internal/models"
	"github.com/wDRxxx/eventflow-backend/internal/repository"
	"github.com/wDRxxx/eventflow-backend/internal/service"
	"github.com/wDRxxx/eventflow-backend/internal/utils"
)

type usersServ struct {
	repo       repository.Repository
	authConfig *config.AuthConfig
}

func NewUsersService(
	repo repository.Repository,
	authConfig *config.AuthConfig,
) service.UsersService {
	s := &usersServ{
		repo:       repo,
		authConfig: authConfig,
	}

	return s
}

func (s *usersServ) RegisterUser(ctx context.Context, user *models.User) error {
	_, err := s.repo.User(ctx, user.Email)
	if err == nil {
		return service.ErrUserAlreadyExists
	}

	if !user.IsOAuth {
		pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
		if err != nil {
			return err
		}
		user.Password = string(pass)
	} else {
		user.Password = " "
	}

	_, err = s.repo.InsertUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *usersServ) Login(ctx context.Context, user *models.User) (string, error) {
	u, err := s.repo.User(ctx, user.Email)
	if err != nil {
		return "", service.ErrWrongCredentials
	}

	if !user.IsOAuth {
		err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
		if err != nil {
			return "", service.ErrWrongCredentials
		}
	}

	accessToken, err := utils.GenerateToken(
		&models.UserClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: fmt.Sprint(u.ID),
			},
			Email: u.Email,
		},
		s.authConfig.RefreshTokenSecret(),
		s.authConfig.RefreshTokenTTL(),
	)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *usersServ) AccessToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := utils.VerifyToken(refreshToken, s.authConfig.RefreshTokenSecret())
	if err != nil {
		return "", err
	}

	accessToken, err := utils.GenerateToken(
		claims,
		s.authConfig.AccessTokenSecret(),
		s.authConfig.AccessTokenTTL(),
	)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *usersServ) User(ctx context.Context, userEmail string) (*models.User, error) {
	user, err := s.repo.User(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *usersServ) UpdateUser(ctx context.Context, user *models.User) error {
	if user.TGUsername != "" {
		err := s.repo.UpdateUserTGUsername(ctx, user.ID, user.TGUsername)
		if err != nil {
			return err
		}
	}

	if user.YookassaSettings.ShopID != "" || user.YookassaSettings.ShopKey != "" {
		user.YookassaSettings.UserID = user.ID
		err := s.repo.UpdateYookassaSettings(ctx, &user.YookassaSettings)
		if err != nil {
			return err
		}
	}

	return nil
}
