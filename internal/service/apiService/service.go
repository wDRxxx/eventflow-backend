package apiService

import (
	"log/slog"
	"sync"

	"github.com/wDRxxx/eventflow-backend/internal/closer"
	"github.com/wDRxxx/eventflow-backend/internal/config"
	"github.com/wDRxxx/eventflow-backend/internal/mailer"
	"github.com/wDRxxx/eventflow-backend/internal/models"
	"github.com/wDRxxx/eventflow-backend/internal/repository"
	"github.com/wDRxxx/eventflow-backend/internal/service"
)

type serv struct {
	wg *sync.WaitGroup

	repo       repository.Repository
	authConfig *config.AuthConfig

	paymentsChan chan *models.TicketPayment
	doneChan     chan struct{}

	mailer mailer.Mailer
}

func NewApiService(
	wg *sync.WaitGroup,
	repo repository.Repository,
	authConfig *config.AuthConfig,
	mailer mailer.Mailer,
) service.ApiService {
	s := &serv{
		wg:           wg,
		repo:         repo,
		authConfig:   authConfig,
		paymentsChan: make(chan *models.TicketPayment),
		doneChan:     make(chan struct{}),
		mailer:       mailer,
	}

	closer.Add(1, func() error {
		slog.Info("sending done signal to api service...")
		s.doneChan <- struct{}{}

		return nil
	})

	closer.Add(2, func() error {
		slog.Info("closing api service channels...")
		close(s.paymentsChan)
		close(s.doneChan)

		return nil
	})

	go s.listenForPayments()

	return s
}
