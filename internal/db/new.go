package db

import (
	"context"
	"fmt"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
	"forge.nhmk.de/khuong/kurz/internal/repository"
)

type UrlRepository interface {
	GetUrl(ctx context.Context, code string) (repository.Shorturl, error)
	CreateUrl(ctx context.Context, params repository.CreateUrlParams) (repository.Shorturl, error)
	ListUrls(ctx context.Context) ([]repository.Shorturl, error)
	DeleteUrl(ctx context.Context, code string) error
	Close()
}

type PostgresRepo struct {
	*repository.Queries
	pool *pgxpool.Pool
}

func (repo *PostgresRepo) Close() {
	repo.pool.Close()
}

func NewRepository(connUrl string) (UrlRepository, error) {
	pool, err := pgxpool.New(context.Background(), connUrl)
	if err != nil {
		return nil, fmt.Errorf("Can't establish a connection pool to the database %w", err)
	}

	m, err := migrate.New(
		"file://db/postgres/migrations",
		connUrl,
	)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("Migration init error: %w", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		pool.Close()
		return nil, fmt.Errorf("Migration error: %w", err)
	}

	return &PostgresRepo {
		Queries: repository.New(pool),
		pool: pool,
	}, nil
}
