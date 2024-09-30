package apiService

import (
	"log"
	"sync"

	"github.com/wDRxxx/yookassa-go-sdk/yookassa"

	"github.com/wDRxxx/eventflow-backend/internal/closer"
	"github.com/wDRxxx/eventflow-backend/internal/config"
	"github.com/wDRxxx/eventflow-backend/internal/models"
	"github.com/wDRxxx/eventflow-backend/internal/repository"
	"github.com/wDRxxx/eventflow-backend/internal/service"
)

type serv struct {
	wg *sync.WaitGroup

	repo       repository.Repository
	authConfig *config.AuthConfig
	yooClient  *yookassa.Client

	paymentsChan chan *models.TicketPayment
	doneChan     chan struct{}
}

func NewApiService(
	wg *sync.WaitGroup,
	repo repository.Repository,
	authConfig *config.AuthConfig,
	yooClient *yookassa.Client,
) service.ApiService {
	s := &serv{
		wg:           wg,
		repo:         repo,
		authConfig:   authConfig,
		yooClient:    yooClient,
		paymentsChan: make(chan *models.TicketPayment),
		doneChan:     make(chan struct{}),
	}

	closer.Add(func() error {
		log.Println("closing payment channels...")
		s.doneChan <- struct{}{}
		close(s.paymentsChan)
		close(s.doneChan)

		return nil
	})

	go s.listenForPayments()

	return s
}
