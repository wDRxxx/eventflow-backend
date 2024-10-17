package httpServer

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/wDRxxx/eventflow-backend/internal/api"
	"github.com/wDRxxx/eventflow-backend/internal/models"
	"github.com/wDRxxx/eventflow-backend/internal/utils"
)

func (s *server) userTickets(w http.ResponseWriter, r *http.Request) {
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

	tickets, err := s.ticketsService.UserTickets(r.Context(), int64(id))
	if err != nil {
		slog.Error("Error getting tickets", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	utils.WriteJSON(tickets, w)
}

func (s *server) buyTicket(w http.ResponseWriter, r *http.Request) {
	_, claims, err := s.getAndVerifyHeaderToken(r)
	if err != nil {
		slog.Error("Error getting claims", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	var req models.BuyTicketRequest
	err = utils.ReadJSON(r.Body, &req)
	if err != nil {
		slog.Error("Error reading request body", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	req.UserEmail = claims.Email

	url, err := s.ticketsService.BuyTicket(r.Context(), &req)
	if err != nil {
		slog.Error("Error buying ticket", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	utils.WriteJSON(&models.DefaultResponse{
		Error:   false,
		Message: url,
	}, w)
}
