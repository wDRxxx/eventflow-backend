package httpServer

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	"github.com/wDRxxx/eventflow-backend/internal/models"
	"github.com/wDRxxx/eventflow-backend/internal/utils"
)

var (
	errInternal = errors.New("Internal error, try again later.")
	errNotFound = errors.New("Not found.")
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

	_, err = s.apiService.CreteEvent(r.Context(), &event)
	if err != nil {
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
	urlTitle := chi.URLParam(r, "url-title")
	event := models.Event{URLTitle: urlTitle}

	err := utils.ReadJSON(w, r, &event)
	if err != nil {
		slog.Error("Error reading request body", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	utils.WriteJSON(&models.DefaultResponse{
		Error:   false,
		Message: "created successfully",
	}, w)
}

func (s *server) deleteEvent(w http.ResponseWriter, r *http.Request) {
	urlTitle := chi.URLParam(r, "url-title")
	err := s.apiService.DeleteEvent(r.Context(), urlTitle)
	if err != nil {
		slog.Error("Error deleting event", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	utils.WriteJSON(&models.DefaultResponse{
		Error:   false,
		Message: "Event was deleted successfully",
	}, w)
}
