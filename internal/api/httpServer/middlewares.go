package httpServer

import (
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/wDRxxx/eventflow-backend/internal/models"
	"github.com/wDRxxx/eventflow-backend/internal/utils"
)

var errInvalidAuthHeader = errors.New("invalid auth header")

func (s *server) authRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Vary", "Authorization")

		_, _, err := s.getAndVerifyHeaderToken(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *server) getAndVerifyHeaderToken(r *http.Request) (string, *models.UserClaims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil, errInvalidAuthHeader
	}

	exploded := strings.Split(authHeader, " ")
	if len(exploded) != 2 || exploded[0] != "Bearer" {
		return "", nil, errInvalidAuthHeader
	}

	token := exploded[1]
	claims, err := utils.VerifyToken(token, s.authConfig.AccessTokenSecret)
	if err != nil {
		return "", nil, errInvalidAuthHeader
	}

	return token, claims, nil
}
