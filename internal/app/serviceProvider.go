package app

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wDRxxx/eventflow-backend/internal/api"
	"github.com/wDRxxx/eventflow-backend/internal/api/httpServer"
	"github.com/wDRxxx/eventflow-backend/internal/config"
	"github.com/wDRxxx/eventflow-backend/internal/repository"
	"github.com/wDRxxx/eventflow-backend/internal/repository/postgres"
	"github.com/wDRxxx/eventflow-backend/internal/service"
	"github.com/wDRxxx/eventflow-backend/internal/service/apiService"
	"log"
)

type serviceProvider struct {
	httpConfig     *config.HttpConfig
	postgresConfig *config.PostgresConfig

	repository repository.Repository
	apiService service.ApiService
	httpServer api.HTTPServer
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) HttpConfig() *config.HttpConfig {
	if s.httpConfig == nil {
		s.httpConfig = config.NewHttpConfig()
	}

	return s.httpConfig
}

func (s *serviceProvider) PostgresConfig() *config.PostgresConfig {
	if s.postgresConfig == nil {
		s.postgresConfig = config.NewPostgresConfig()
	}
	return s.postgresConfig
}

func (s *serviceProvider) Repository(ctx context.Context) repository.Repository {
	if s.repository == nil {
		db, err := pgxpool.New(ctx, s.PostgresConfig().ConnectionString())
		if err != nil {
			log.Fatalf("error connecting to database: %v", err)
		}
		err = db.Ping(ctx)
		if err != nil {
			log.Fatalf("error connecting to database: %v", err)
		}

		s.repository = postgres.NewPostgresRepo(db)
	}

	return s.repository
}

func (s *serviceProvider) ApiService(ctx context.Context) service.ApiService {
	if s.apiService == nil {
		s.apiService = apiService.NewApiService(s.Repository(ctx))
	}

	return s.apiService
}

func (s *serviceProvider) HTTPServer(ctx context.Context) api.HTTPServer {
	if s.httpServer == nil {
		s.httpServer = httpServer.NewHTTPServer(s.ApiService(ctx))
	}

	return s.httpServer
}