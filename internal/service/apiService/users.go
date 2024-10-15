package apiService

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/wDRxxx/eventflow-backend/internal/models"
	"github.com/wDRxxx/eventflow-backend/internal/service"
	"github.com/wDRxxx/eventflow-backend/internal/utils"
)

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

func (s *serv) AccessToken(ctx context.Context, refreshToken string) (string, error) {
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

func (s *serv) User(ctx context.Context, userEmail string) (*models.User, error) {
	user, err := s.repo.User(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *serv) UpdateUser(ctx context.Context, user *models.User) error {
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
