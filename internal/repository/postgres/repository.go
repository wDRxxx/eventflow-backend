package postgres

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/wDRxxx/eventflow-backend/internal/repository"
)

type repo struct {
	db      *pgxpool.Pool
	timeout time.Duration
}

const (
	usersTable            = "users"
	eventsTable           = "events"
	pricesTable           = "prices"
	ticketsTable          = "tickets"
	yookassaSettingsTable = "users_yookassa_settings"

	structTag = "db"
)

func NewPostgresRepo(db *pgxpool.Pool, timeout time.Duration) repository.Repository {
	return &repo{
		db:      db,
		timeout: timeout,
	}
}
