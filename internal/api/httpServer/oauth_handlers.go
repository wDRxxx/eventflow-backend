package httpServer

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"github.com/wDRxxx/eventflow-backend/internal/models"
	"github.com/wDRxxx/eventflow-backend/internal/service"
	"github.com/wDRxxx/eventflow-backend/internal/utils"
)

const (
	googleUserInfoEndpoint = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
)

func (s *server) oauthCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	oauthStateCookie, _ := r.Cookie("oauth_state")

	state := r.URL.Query().Get("state")
	if state != oauthStateCookie.Value {
		utils.WriteJSONError(errors.New("wrong state"), w, http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	user := &models.User{
		IsOAuth: true,
	}

	switch provider {
	case "google":
		token, err := s.oauth.GoogleConfig().Exchange(r.Context(), code)
		if err != nil {
			slog.Error("error while exchanging", slog.Any("error", err))
			utils.WriteJSONError(err, w, http.StatusBadRequest)
			return
		}

		resp, err := http.Get(googleUserInfoEndpoint + token.AccessToken)
		if err != nil {
			slog.Error("error while getting response from google", slog.Any("error", err))
			utils.WriteJSONError(err, w, http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()

		userData, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("error reading response body", slog.Any("error", err))
			utils.WriteJSONError(err, w, http.StatusBadRequest)
		}

		err = json.Unmarshal(userData, &user)
		if err != nil {
			slog.Error("error parsing response", slog.Any("error", err))
			utils.WriteJSONError(err, w, http.StatusBadRequest)
		}
	}

	err := s.usersService.RegisterUser(r.Context(), user)
	if err != nil {
		if !errors.Is(err, service.ErrUserAlreadyExists) {
			slog.Error("error registering user", slog.Any("error", err))
			return
		}
	}

	refreshToken, err := s.usersService.Login(r.Context(), user)
	if err != nil {
		slog.Error("error logging in", slog.Any("error", err))
		utils.WriteJSONError(err, w)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(s.authConfig.RefreshTokenTTL()),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, s.oauth.RedirectURL(), http.StatusFound)
}

func (s *server) oauthLogin(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	var url string

	var expiration = time.Now().Add(20 * time.Minute)
	b := make([]byte, 16)
	rand.Read(b)

	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{
		Name:    "oauth_state",
		Value:   state,
		Expires: expiration,
	}

	switch provider {
	case "google":
		url = s.oauth.GoogleConfig().AuthCodeURL(state)
		if strings.Index(url, "https") == -1 {
			url = strings.Replace(url, "http", "https", 1)
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, url, http.StatusSeeOther)
}
