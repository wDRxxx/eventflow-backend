package ticketsService

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/wDRxxx/yookassa-go-sdk/yookassa"
	yoomodels "github.com/wDRxxx/yookassa-go-sdk/yookassa/models"
	yoopayment "github.com/wDRxxx/yookassa-go-sdk/yookassa/models/payment"

	"github.com/wDRxxx/eventflow-backend/internal/closer"
	"github.com/wDRxxx/eventflow-backend/internal/config"
	"github.com/wDRxxx/eventflow-backend/internal/mailer"
	"github.com/wDRxxx/eventflow-backend/internal/models"
	"github.com/wDRxxx/eventflow-backend/internal/repository"
	"github.com/wDRxxx/eventflow-backend/internal/service"
)

type ticketsServ struct {
	wg *sync.WaitGroup

	repo       repository.Repository
	authConfig *config.AuthConfig

	paymentsChan chan *models.TicketPayment
	doneChan     chan struct{}

	mailer mailer.Mailer
}

func NewTicketsService(
	wg *sync.WaitGroup,
	repo repository.Repository,
	mailer mailer.Mailer,
) service.TicketsService {
	s := &ticketsServ{
		wg:           wg,
		repo:         repo,
		paymentsChan: make(chan *models.TicketPayment),
		doneChan:     make(chan struct{}),
		mailer:       mailer,
	}

	closer.Add(1, func() error {
		slog.Info("sending done signal to tickets service...")
		s.doneChan <- struct{}{}

		return nil
	})

	closer.Add(2, func() error {
		slog.Info("closing tickets service channels...")
		close(s.paymentsChan)
		close(s.doneChan)

		return nil
	})

	go s.listenForPayments()

	return s
}

func (s *ticketsServ) listenForPayments() {
	for {
		select {
		case ticketPayment := <-s.paymentsChan:
			s.wg.Add(1)

			go func() {
				s.wg.Done()

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

func (s *ticketsServ) checkPaymentStatus(ticketPayment *models.TicketPayment) error {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	for {
		select {
		case <-ticker.C:
			payment, err := ticketPayment.YooClient.PaymentInfo(ticketPayment.Payment.ID)
			if err != nil {
				return err
			}

			if payment.Status == "succeeded" {
				ticket := &models.Ticket{
					ID:        uuid.NewString(),
					UserID:    ticketPayment.User.ID,
					User:      ticketPayment.User,
					EventID:   ticketPayment.Event.ID,
					Event:     ticketPayment.Event,
					IsUsed:    false,
					FirstName: ticketPayment.BuyTicketRequest.FirstName,
					LastName:  ticketPayment.BuyTicketRequest.LastName,
					PaymentID: payment.ID,
				}

				err = s.createTicket(ctx, ticket)
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

func (s *ticketsServ) BuyTicket(ctx context.Context, req *models.BuyTicketRequest) (string, error) {
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
			User:      user,
			EventID:   event.ID,
			Event:     event,
			IsUsed:    false,
			FirstName: req.FirstName,
			LastName:  req.LastName,
		}

		err = s.createTicket(ctx, ticket)
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

	yooClient := yookassa.NewClient(event.ShopID, event.ShopKey)
	payment := &yoopayment.Payment{
		Amount: amount,
		PaymentMethodData: &yoopayment.PaymentMethodData{
			Type: "bank_card",
		},
		Confirmation: &yoomodels.Confirmation{
			Type:      "redirect",
			ReturnURL: s.authConfig.Domain() + "/user/profile",
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

	respPayment, err := yooClient.CreatePayment(payment)
	if err != nil {
		return "", err
	}

	s.paymentsChan <- &models.TicketPayment{
		Payment:          respPayment,
		BuyTicketRequest: req,
		User:             user,
		Event:            event,
		Ctx:              ctx,
		YooClient:        yooClient,
	}

	return respPayment.Confirmation.ConfirmationURL, nil
}

func (s *ticketsServ) createTicket(ctx context.Context, ticket *models.Ticket) error {
	id, err := s.repo.InsertTicket(ctx, ticket)
	if err != nil {
		return err
	}

	msg := &models.OrderMessage{
		To:          []string{ticket.User.Email},
		TicketID:    id,
		EventTitle:  ticket.Event.Title,
		ImageURL:    s.authConfig.Domain() + "/api/static/" + ticket.Event.PreviewImage,
		RedirectURL: s.authConfig.Domain() + "/user/profile",
	}

	s.mailer.SendOrderMail(msg)

	return nil
}

func (s *ticketsServ) Ticket(ctx context.Context, ticketID string) (*models.Ticket, error) {
	ticket, err := s.repo.Ticket(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (s *ticketsServ) UserTickets(ctx context.Context, userID int64) ([]*models.Ticket, error) {
	tickets, err := s.repo.UserTickets(ctx, userID)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}
