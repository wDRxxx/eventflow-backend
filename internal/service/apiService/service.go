package apiService

import (
	"github.com/wDRxxx/eventflow-backend/internal/repository"
	"github.com/wDRxxx/eventflow-backend/internal/service"
)

type serv struct {
	repo repository.Repository
}

func NewApiService(repo repository.Repository) service.ApiService {
	return &serv{
		repo: repo,
	}
}
