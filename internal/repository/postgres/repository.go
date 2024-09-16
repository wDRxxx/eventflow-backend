package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wDRxxx/eventflow-backend/internal/repository"
)

type repo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(db *pgxpool.Pool) repository.Repository {
	return &repo{
		db: db,
	}
}
