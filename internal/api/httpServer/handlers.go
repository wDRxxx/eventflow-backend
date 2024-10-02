package httpServer

import (
	"bytes"
	"log/slog"
	"net/http"
	"strconv"
	"time"

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
			utils.WriteJSONError(errInternal, w)
			return
		}
	}

	utils.WriteJSON(events, w)
}

func (s *server) myEvents(w http.ResponseWriter, r *http.Request) {
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

	events, err := s.apiService.UserEvents(r.Context(), id)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			slog.Error("Error getting events", slog.Any("error", err))
			utils.WriteJSONError(errInternal, w)
			return
		}
	}

	utils.WriteJSON(events, w)
}

func (s *server) createEvent(w http.ResponseWriter, r *http.Request) {
	_, claims, err := s.getAndVerifyHeaderToken(r)
	if err != nil {
		slog.Error("Error getting claims", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		slog.Error("Error reading request body", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	imgs, err := saveMultipartImages(r, "image", s.httpConfig.StaticDir())
	if err != nil {
		slog.Error("Error saving event image", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	var event models.Event
	err = utils.ReadJSON(bytes.NewBuffer([]byte(r.Form.Get("event"))), &event)
	if err != nil {
		slog.Error("Error reading request body", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	if event.Title == "" ||
		event.Description == "" ||
		event.BeginningTime.IsZero() ||
		event.EndTime.IsZero() ||
		event.Location == "" {

		utils.WriteJSONError(errWrongInput, w, http.StatusUnprocessableEntity)
		return
	}

	event.PreviewImage = imgs[0]

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

	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		slog.Error("Error reading request body", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}
	imgs, err := saveMultipartImages(r, "image", s.httpConfig.StaticDir())
	if err != nil {
		slog.Error("Error saving event image", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	event := models.Event{URLTitle: urlTitle}
	err = utils.ReadJSON(bytes.NewBuffer([]byte(r.Form.Get("event"))), &event)
	if err != nil {
		slog.Error("Error reading request body", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	if event.Title == "" ||
		event.Description == "" ||
		event.BeginningTime.IsZero() ||
		event.EndTime.IsZero() ||
		event.Location == "" {

		utils.WriteJSONError(errWrongInput, w, http.StatusUnprocessableEntity)
		return
	}

	if len(imgs) > 0 {
		event.PreviewImage = imgs[0]
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
		utils.WriteJSONError(errors.New("Error. Maybe someone already has bought a ticket..."), w)
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

func (s *server) userTickets(w http.ResponseWriter, r *http.Request) {
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

	tickets, err := s.apiService.UserTickets(r.Context(), int64(id))
	if err != nil {
		slog.Error("Error getting tickets", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	utils.WriteJSON(tickets, w)
}

func (s *server) buyTicket(w http.ResponseWriter, r *http.Request) {
	_, claims, err := s.getAndVerifyHeaderToken(r)
	if err != nil {
		slog.Error("Error getting claims", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	var req models.BuyTicketRequest
	err = utils.ReadJSON(r.Body, &req)
	if err != nil {
		slog.Error("Error reading request body", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	req.UserEmail = claims.Email

	url, err := s.apiService.BuyTicket(r.Context(), &req)
	if err != nil {
		slog.Error("Error buying ticket", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	utils.WriteJSON(&models.DefaultResponse{
		Error:   false,
		Message: url,
	}, w)
}

//

func (s *server) register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := utils.ReadReqJSON(w, r, &user)
	if err != nil {
		slog.Error("Error reading request body", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	if !utils.IsEmail(user.Email) ||
		user.Password == "" {

		utils.WriteJSONError(errWrongInput, w, http.StatusUnprocessableEntity)
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
	err := utils.ReadReqJSON(w, r, &user)
	if err != nil {
		slog.Error("Error reading request body", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	if !utils.IsEmail(user.Email) ||
		user.Password == "" {

		utils.WriteJSONError(errWrongInput, w, http.StatusUnprocessableEntity)
		return
	}

	refreshToken, err := s.apiService.Login(r.Context(), &user)
	if err != nil {
		if errors.Is(err, service.ErrWrongCredentials) {
			utils.WriteJSONError(err, w, http.StatusUnauthorized)
			return
		}

		slog.Error("Error getting token", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}
	accessToken, err := s.apiService.AccessToken(r.Context(), refreshToken)

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
		utils.WriteJSONError(errInternal, w)
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
		utils.WriteJSONError(errInternal, w)
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
		utils.WriteJSONError(errInternal, w)
		return
	}

	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		slog.Error("Error parsing subject", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}
	user.ID = int64(id)

	err = s.apiService.UpdateUser(r.Context(), &user)
	if err != nil {
		slog.Error("Error updating user", slog.Any("error", err))
		utils.WriteJSONError(errInternal, w)
		return
	}

	utils.WriteJSON(&models.DefaultResponse{
		Error:   false,
		Message: "Successfully updated user",
	}, w)
}

func saveMultipartImages(r *http.Request, formField string, staticDir string) ([]string, error) {
	reqImages := r.MultipartForm.File[formField]
	var images []string

	for _, img := range reqImages {
		filename, err := utils.SaveStaticImage(img, staticDir)
		if err != nil {
			return nil, err
		}

		images = append(images, filename)
	}

	return images, nil
}
