package apiService

import (
	"github.com/wDRxxx/yookassa-go-sdk/yookassa"

	"github.com/wDRxxx/eventflow-backend/internal/config"
	"github.com/wDRxxx/eventflow-backend/internal/repository"
	"github.com/wDRxxx/eventflow-backend/internal/service"
)

type serv struct {
	repo       repository.Repository
	authConfig *config.AuthConfig
	yooClient  *yookassa.Client

	paymentsChan chan string
	doneChan     chan struct{}
}

func NewApiService(
	repo repository.Repository,
	authConfig *config.AuthConfig,
	yooClient *yookassa.Client,
) service.ApiService {
	return &serv{
		repo:         repo,
		authConfig:   authConfig,
		yooClient:    yooClient,
		paymentsChan: make(chan string),
		doneChan:     make(chan struct{}),
	}
}
