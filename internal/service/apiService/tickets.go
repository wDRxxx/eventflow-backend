package apiService

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	yoomodels "github.com/wDRxxx/yookassa-go-sdk/yookassa/models"
	yoopayment "github.com/wDRxxx/yookassa-go-sdk/yookassa/models/payment"

	"github.com/wDRxxx/eventflow-backend/internal/models"
)

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

		_, err = s.repo.InsertTicket(ctx, ticket)
		if err != nil {
			return "", err
		}

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

	// generate payment
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
	}

	respPayment, err := s.yooClient.CreatePayment(payment)
	if err != nil {
		return "", err
	}
	// listen payment in payment chan

	// return payment url
	return respPayment.Confirmation.ConfirmationURL, nil
}

func (s *serv) Ticket(ctx context.Context, ticketID string) (*models.Ticket, error) {
	ticket, err := s.repo.Ticket(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}
