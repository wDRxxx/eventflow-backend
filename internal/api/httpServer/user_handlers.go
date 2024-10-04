package httpServer

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/wDRxxx/eventflow-backend/internal/api"
	"github.com/wDRxxx/eventflow-backend/internal/models"
	"github.com/wDRxxx/eventflow-backend/internal/service"
	"github.com/wDRxxx/eventflow-backend/internal/utils"
)

func (s *server) register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := utils.ReadReqJSON(w, r, &user)
	if err != nil {
		slog.Error("Error reading request body", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	if !utils.IsEmail(user.Email) ||
		user.Password == "" {

		utils.WriteJSONError(api.ErrWrongInput, w, http.StatusUnprocessableEntity)
		return
	}

	err = s.apiService.RegisterUser(r.Context(), &user)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			utils.WriteJSONError(err, w, http.StatusConflict)
			return
		}

		slog.Error("Error registering user", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	utils.WriteJSON(&models.DefaultResponse{
		Error:   false,
		Message: "You was successfully registered",
	}, w, http.StatusCreated)
}

func (s *server) login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := utils.ReadReqJSON(w, r, &user)
	if err != nil {
		slog.Error("Error reading request body", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	if !utils.IsEmail(user.Email) ||
		user.Password == "" {

		utils.WriteJSONError(api.ErrWrongInput, w, http.StatusUnprocessableEntity)
		return
	}

	refreshToken, err := s.apiService.Login(r.Context(), &user)
	if err != nil {
		if errors.Is(err, service.ErrWrongCredentials) {
			utils.WriteJSONError(err, w, http.StatusUnauthorized)
			return
		}

		slog.Error("Error getting token", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}
	accessToken, err := s.apiService.AccessToken(r.Context(), refreshToken)
	if err != nil {
		slog.Error("Error getting access token", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
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

	utils.WriteJSON(&models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, w)
}

func (s *server) refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	accessToken, err := s.apiService.AccessToken(r.Context(), refreshToken.Value)
	if err != nil {
		slog.Error("Error getting access token", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	utils.WriteJSON(&models.DefaultResponse{
		Error:   false,
		Message: accessToken,
	}, w)
}

func (s *server) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusAccepted)
}

func (s *server) profile(w http.ResponseWriter, r *http.Request) {
	_, claims, err := s.getAndVerifyHeaderToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := s.apiService.User(r.Context(), claims.Email)
	if err != nil {
		slog.Error("Error getting user", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}
	user.Password = ""

	utils.WriteJSON(user, w)
}

func (s *server) updateProfile(w http.ResponseWriter, r *http.Request) {
	_, claims, err := s.getAndVerifyHeaderToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var user models.User
	err = utils.ReadReqJSON(w, r, &user)
	if err != nil {
		slog.Error("Error reading request body", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		slog.Error("Error parsing subject", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}
	user.ID = int64(id)

	err = s.apiService.UpdateUser(r.Context(), &user)
	if err != nil {
		slog.Error("Error updating user", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	utils.WriteJSON(&models.DefaultResponse{
		Error:   false,
		Message: "Successfully updated user",
	}, w)
}
