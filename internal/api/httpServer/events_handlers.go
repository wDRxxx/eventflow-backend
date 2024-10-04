package httpServer

import (
	"bytes"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	"github.com/wDRxxx/eventflow-backend/internal/api"
	"github.com/wDRxxx/eventflow-backend/internal/models"
	"github.com/wDRxxx/eventflow-backend/internal/service"
	"github.com/wDRxxx/eventflow-backend/internal/utils"
)

func (s *server) event(w http.ResponseWriter, r *http.Request) {
	urlTitle := chi.URLParam(r, "url-title")
	resp, err := s.apiService.Event(r.Context(), urlTitle)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			utils.WriteJSONError(api.ErrNotFound, w, http.StatusNotFound)
			return
		}

		slog.Error("Error getting event", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	utils.WriteJSON(resp, w)
}

func (s *server) events(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Query().Get("page")
	page, err := strconv.Atoi(p)
	if err != nil || page < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	events, err := s.apiService.Events(r.Context(), page)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			slog.Error("Error getting events", slog.Any("error", err))
			utils.WriteJSONError(api.ErrInternal, w)
			return
		}
	}

	utils.WriteJSON(events, w)
}

func (s *server) myEvents(w http.ResponseWriter, r *http.Request) {
	_, claims, err := s.getAndVerifyHeaderToken(r)
	if err != nil {
		slog.Error("Error getting claims", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}
	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		slog.Error("Error converting claims.Subject to int", slog.Any("error", err), slog.String("subject", claims.Subject))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	events, err := s.apiService.UserEvents(r.Context(), int64(id))
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			slog.Error("Error getting events", slog.Any("error", err))
			utils.WriteJSONError(api.ErrInternal, w)
			return
		}
	}

	utils.WriteJSON(events, w)
}

func (s *server) createEvent(w http.ResponseWriter, r *http.Request) {
	_, claims, err := s.getAndVerifyHeaderToken(r)
	if err != nil {
		slog.Error("Error getting claims", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		slog.Error("Error reading request body", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	imgs, err := saveMultipartImages(r, "image", s.httpConfig.StaticDir())
	if err != nil {
		slog.Error("Error saving event image", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	var event models.Event
	err = utils.ReadJSON(bytes.NewBuffer([]byte(r.Form.Get("event"))), &event)
	if err != nil {
		slog.Error("Error reading request body", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	if event.Title == "" ||
		event.Description == "" ||
		event.BeginningTime.IsZero() ||
		event.EndTime.IsZero() ||
		event.Location == "" {

		utils.WriteJSONError(api.ErrWrongInput, w, http.StatusUnprocessableEntity)
		return
	}

	if len(imgs) > 0 {
		event.PreviewImage = imgs[0]
	}

	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		slog.Error("Error converting claims.Subject to int", slog.Any("error", err), slog.String("subject", claims.Subject))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	event.CreatorID = int64(id)

	_, err = s.apiService.CreateEvent(r.Context(), &event)
	if err != nil {
		if errors.Is(err, service.ErrNoPrices) || errors.Is(err, service.ErrPricesForFree) {
			utils.WriteJSONError(err, w)
			return
		}

		slog.Error("Error creating new event", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
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
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}
	id, err := strconv.Atoi(claims.Subject)
	urlTitle := chi.URLParam(r, "url-title")

	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		slog.Error("Error reading request body", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}
	imgs, err := saveMultipartImages(r, "image", s.httpConfig.StaticDir())
	if err != nil {
		slog.Error("Error saving event image", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	event := models.Event{URLTitle: urlTitle}
	err = utils.ReadJSON(bytes.NewBuffer([]byte(r.Form.Get("event"))), &event)
	if err != nil {
		slog.Error("Error reading request body", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	if len(imgs) > 0 {
		event.PreviewImage = imgs[0]
	}

	err = s.apiService.UpdateEvent(r.Context(), int64(id), &event)
	if err != nil {
		slog.Error("Error updating event", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
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
		utils.WriteJSONError(api.ErrInternal, w)
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
		utils.WriteJSONError(errors.New("Error. Maybe someone already has bought a ticket..."), w)
		return
	}

	utils.WriteJSON(&models.DefaultResponse{
		Error:   false,
		Message: "Event was deleted successfully",
	}, w)
}
