package apiService

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	yoomodels "github.com/wDRxxx/yookassa-go-sdk/yookassa/models"
	yoopayment "github.com/wDRxxx/yookassa-go-sdk/yookassa/models/payment"

	"github.com/wDRxxx/eventflow-backend/internal/models"
	"github.com/wDRxxx/eventflow-backend/internal/service"
)

func (s *serv) listenForPayments() {
	for {
		select {
		case ticketPayment := <-s.paymentsChan:
			s.wg.Add(1)

			go func() {
				err := s.checkPaymentStatus(ticketPayment)
				if err != nil {
					if !errors.Is(err, service.ErrPaymentTimeout) {
						slog.Error("error checking payment", slog.Any("error", err))
					}
				}
			}()
		case <-s.doneChan:
			return
		}
	}
}

func (s *serv) checkPaymentStatus(ticketPayment *models.TicketPayment) error {
	defer s.wg.Done()

	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	for {
		select {
		case <-ticker.C:
			payment, err := s.yooClient.PaymentInfo(ticketPayment.Payment.ID)
			if err != nil {
				return err
			}

			if payment.Status == "succeeded" {
				ticket := &models.Ticket{
					ID:        uuid.NewString(),
					UserID:    ticketPayment.User.ID,
					EventID:   ticketPayment.Event.ID,
					IsUsed:    false,
					FirstName: ticketPayment.BuyTicketRequest.FirstName,
					LastName:  ticketPayment.BuyTicketRequest.LastName,
					PaymentID: payment.ID,
				}

				err := s.CreateTicket(ctx, ticket)
				if err != nil {
					return err
				}

				return nil
			}
		case <-ctx.Done():
			return service.ErrPaymentTimeout
		}
	}
}

func (s *serv) BuyTicket(ctx context.Context, req *models.BuyTicketRequest) (string, error) {
	user, err := s.repo.User(ctx, req.UserEmail)
	if err != nil {
		return "", err
	}

	event, err := s.repo.EventByURLTitle(ctx, req.EventUrlTitle)
	if err != nil {
		return "", err
	}

	if event.IsFree {
		ticket := &models.Ticket{
			ID:        uuid.NewString(),
			UserID:    user.ID,
			EventID:   event.ID,
			IsUsed:    false,
			FirstName: req.FirstName,
			LastName:  req.LastName,
		}

		err = s.CreateTicket(ctx, ticket)
		return "", nil
	}
	var amount *yoomodels.Amount
	for _, p := range event.Prices {
		if p.ID == req.PriceID {
			amount = &yoomodels.Amount{
				Value:    fmt.Sprintf("%d.00", p.Price),
				Currency: p.Currency,
			}
			break
		}
	}

	payment := &yoopayment.Payment{
		Amount: amount,
		PaymentMethodData: &yoopayment.PaymentMethodData{
			Type: "bank_card",
		},
		// TODO: change return url
		Confirmation: &yoomodels.Confirmation{
			Type:      "redirect",
			ReturnURL: "localhost:3000/user/profile",
		},
		Receipt: &yoomodels.Receipt{
			Email: user.Email,
			Items: []*yoomodels.Item{
				{
					Description: fmt.Sprintf("Ticket to \"%s\" event", event.Title),
					Amount:      amount,
					Quantity:    1,
					VATCode:     1,
				},
			},
		},
		Capture: true,
	}

	respPayment, err := s.yooClient.CreatePayment(payment)
	if err != nil {
		return "", err
	}

	s.paymentsChan <- &models.TicketPayment{
		Payment:          respPayment,
		BuyTicketRequest: req,
		User:             user,
		Event:            event,
		Ctx:              ctx,
	}

	return respPayment.Confirmation.ConfirmationURL, nil
}

func (s *serv) CreateTicket(ctx context.Context, ticket *models.Ticket) error {
	_, err := s.repo.InsertTicket(ctx, ticket)
	if err != nil {
		return err
	}

	return nil
}

func (s *serv) Ticket(ctx context.Context, ticketID string) (*models.Ticket, error) {
	ticket, err := s.repo.Ticket(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}