package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/wDRxxx/eventflow-backend/internal/api"
	"github.com/wDRxxx/eventflow-backend/internal/api/httpServer"
	"github.com/wDRxxx/eventflow-backend/internal/closer"
	"github.com/wDRxxx/eventflow-backend/internal/config"
	"github.com/wDRxxx/eventflow-backend/internal/mailer"
	"github.com/wDRxxx/eventflow-backend/internal/mailer/smtp"
	"github.com/wDRxxx/eventflow-backend/internal/repository"
	"github.com/wDRxxx/eventflow-backend/internal/repository/postgres"
	"github.com/wDRxxx/eventflow-backend/internal/service"
	"github.com/wDRxxx/eventflow-backend/internal/service/apiService"
)

type serviceProvider struct {
	httpConfig     *config.HttpConfig
	postgresConfig *config.PostgresConfig
	authConfig     *config.AuthConfig
	mailerConfig   *config.MailerConfig
	metricsConfig  *config.MetricsConfig

	repository repository.Repository
	apiService service.ApiService
	httpServer api.HTTPServer

	mailer mailer.Mailer
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

func (s *serviceProvider) AuthConfig() *config.AuthConfig {
	if s.authConfig == nil {
		s.authConfig = config.NewAuthConfig()
	}

	return s.authConfig
}

func (s *serviceProvider) MailerConfig() *config.MailerConfig {
	if s.mailerConfig == nil {
		s.mailerConfig = config.NewMailerConfig()
	}

	return s.mailerConfig
}

func (s *serviceProvider) MetricsConfig() *config.MetricsConfig {
	if s.metricsConfig == nil {
		s.metricsConfig = config.NewMetricsConfig()
	}

	return s.metricsConfig
}

func (s *serviceProvider) Repository(ctx context.Context) repository.Repository {
	if s.repository == nil {
		db, err := pgxpool.New(ctx, s.PostgresConfig().ConnectionString())
		if err != nil {
			log.Fatalf("error connecting to database: %v", err)
		}
		closer.Add(2, func() error {
			slog.Info("closing pgxpool")
			db.Close()
			return nil
		})

		err = db.Ping(ctx)
		if err != nil {
			log.Fatalf("error connecting to database: %v", err)
		}

		s.repository = postgres.NewPostgresRepo(db, s.PostgresConfig().Timeout())
	}

	return s.repository
}

func (s *serviceProvider) ApiService(ctx context.Context, wg *sync.WaitGroup) service.ApiService {
	if s.apiService == nil {
		s.apiService = apiService.NewApiService(
			wg,
			s.Repository(ctx),
			s.AuthConfig(),
			s.Mailer(wg),
		)
	}

	return s.apiService
}

func (s *serviceProvider) HTTPServer(ctx context.Context, wg *sync.WaitGroup) api.HTTPServer {
	if s.httpServer == nil {
		s.httpServer = httpServer.NewHTTPServer(s.ApiService(ctx, wg), s.AuthConfig(), s.HttpConfig())
	}

	return s.httpServer
}

func (s *serviceProvider) Mailer(wg *sync.WaitGroup) mailer.Mailer {
	if s.mailer == nil {
		m, err := smtp.NewSMTPMailer(s.MailerConfig(), wg)
		if err != nil {
			panic(fmt.Sprintf("error creating mailer: %v", err))
		}

		s.mailer = m
	}

	return s.mailer
}
