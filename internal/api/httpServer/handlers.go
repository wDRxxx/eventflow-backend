package httpServer

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	"github.com/wDRxxx/eventflow-backend/internal/models"
	"github.com/wDRxxx/eventflow-backend/internal/service"
	"github.com/wDRxxx/eventflow-backend/internal/utils"
)

func (s *server) home(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello World"))
	if err != nil {
		slog.Error("Error writing response", slog.Any("error", err))
	}
}

func (s *server) event(w http.ResponseWriter, r *http.Request) {
	urlTitle := chi.URLParam(r, "url-title")
	resp, err := s.apiService.Event(r.Context(), urlTitle)

	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			slog.Error("Error getting event", slog.Any("error", err))
			utils.WriteJSONError(errInternal, w)
			return
		}

		utils.WriteJSONError(errNotFound, w)
		return
	}

	utils.WriteJSON(resp, w)
}

func (s *server) createEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	err := utils.ReadJSON(w, r, &event)
	if err != nil {
		slog.Error("Error reading request body", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	_, claims, err := s.getAndVerifyHeaderToken(r)
	if err != nil {
		slog.Error("Error getting claims", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		slog.Error("Error converting claims.Subject to int", slog.Any("error", err), slog.String("subject", claims.Subject))
		utils.WriteJSONError(errInternal, w)
		return
	}

	event.CreatorID = int64(id)

	_, err = s.apiService.CreteEvent(r.Context(), &event)
	if err != nil {
		if errors.Is(err, service.ErrNoPrices) || errors.Is(err, service.ErrPricesForFree) {
			utils.WriteJSONError(err, w)
			return
		}

		slog.Error("Error creating new event", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	utils.WriteJSON(&models.DefaultResponse{
		Error:   false,
		Message: "created successfully",
	}, w, http.StatusCreated)
}

func (s *server) updateEvent(w http.ResponseWriter, r *http.Request) {
	_, claims, err := s.getAndVerifyHeaderToken(r)
	if err != nil {
		slog.Error("Error getting claims", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}
	id, err := strconv.Atoi(claims.Subject)
	urlTitle := chi.URLParam(r, "url-title")

	event := models.Event{URLTitle: urlTitle}
	err = utils.ReadJSON(w, r, &event)
	if err != nil {
		slog.Error("Error reading request body", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	err = s.apiService.UpdateEvent(r.Context(), int64(id), &event)
	if err != nil {
		slog.Error("Error updating event", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	utils.WriteJSON(&models.DefaultResponse{
		Error:   false,
		Message: "updated successfully",
	}, w)
}

func (s *server) deleteEvent(w http.ResponseWriter, r *http.Request) {
	_, claims, err := s.getAndVerifyHeaderToken(r)
	if err != nil {
		slog.Error("Error getting claims", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}
	id, err := strconv.Atoi(claims.Subject)
	urlTitle := chi.URLParam(r, "url-title")

	err = s.apiService.DeleteEvent(r.Context(), int64(id), urlTitle)
	if err != nil {
		if errors.Is(err, service.ErrPermissionDenied) {
			utils.WriteJSONError(service.ErrPermissionDenied, w, http.StatusForbidden)
			return
		}

		slog.Error("Error deleting event", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	utils.WriteJSON(&models.DefaultResponse{
		Error:   false,
		Message: "Event was deleted successfully",
	}, w)
}

//

func (s *server) ticket(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ticket, err := s.apiService.Ticket(r.Context(), id)
	if err != nil {
		slog.Error("Error getting ticket", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	utils.WriteJSON(ticket, w)

}

//

func (s *server) register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := utils.ReadJSON(w, r, &user)
	if err != nil {
		slog.Error("Error reading request body", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	err = s.apiService.RegisterUser(r.Context(), &user)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			utils.WriteJSONError(err, w, http.StatusConflict)
			return
		}

		slog.Error("Error registering user", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	utils.WriteJSON(&models.DefaultResponse{
		Error:   false,
		Message: "You was successfully registered",
	}, w, http.StatusCreated)
}

func (s *server) login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := utils.ReadJSON(w, r, &user)
	if err != nil {
		slog.Error("Error reading request body", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	token, err := s.apiService.Login(r.Context(), &user)
	if err != nil {
		if errors.Is(err, service.ErrWrongCredentials) {
			utils.WriteJSONError(err, w, http.StatusUnauthorized)
			return
		}

		slog.Error("Error getting token", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	utils.WriteJSON(&models.DefaultResponse{
		Error:   false,
		Message: token,
	}, w)
}

func (s *server) refresh(w http.ResponseWriter, r *http.Request) {
	var t struct {
		RefreshToken string `json:"refresh_token"`
	}

	err := utils.ReadJSON(w, r, &t)
	if err != nil {
		slog.Error("Error reading request body", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	accessToken, err := s.apiService.AccessToken(r.Context(), t.RefreshToken)
	if err != nil {
		slog.Error("Error getting access token", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	utils.WriteJSON(&models.DefaultResponse{
		Error:   false,
		Message: accessToken,
	}, w)
}
