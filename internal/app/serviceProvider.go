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
	"github.com/wDRxxx/eventflow-backend/internal/service/eventsService"
	"github.com/wDRxxx/eventflow-backend/internal/service/ticketsService"
	"github.com/wDRxxx/eventflow-backend/internal/service/usersService"
)

type serviceProvider struct {
	httpConfig     *config.HttpConfig
	postgresConfig *config.PostgresConfig
	authConfig     *config.AuthConfig
	mailerConfig   *config.MailerConfig
	metricsConfig  *config.MetricsConfig

	repository repository.Repository
	httpServer api.HTTPServer

	eventsService  service.EventsService
	ticketsService service.TicketsService
	usersService   service.UsersService

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

func (s *serviceProvider) EventsService(ctx context.Context) service.EventsService {
	if s.eventsService == nil {
		s.eventsService = eventsService.NewEventsService(s.Repository(ctx))
	}

	return s.eventsService
}

func (s *serviceProvider) TicketsService(ctx context.Context, wg *sync.WaitGroup) service.TicketsService {
	if s.ticketsService == nil {
		s.ticketsService = ticketsService.NewTicketsService(
			wg,
			s.Repository(ctx),
			s.Mailer(wg),
		)
	}

	return s.ticketsService
}

func (s *serviceProvider) UsersService(ctx context.Context) service.UsersService {
	if s.usersService == nil {
		s.usersService = usersService.NewUsersService(s.Repository(ctx), s.AuthConfig())
	}

	return s.usersService
}

func (s *serviceProvider) HTTPServer(ctx context.Context, wg *sync.WaitGroup) api.HTTPServer {
	if s.httpServer == nil {
		s.httpServer = httpServer.NewHTTPServer(
			s.AuthConfig(),
			s.HttpConfig(),
			s.EventsService(ctx),
			s.TicketsService(ctx, wg),
			s.UsersService(ctx),
		)
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
