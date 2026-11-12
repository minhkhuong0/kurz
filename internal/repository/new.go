package repository

import (
	"context"
	"fmt"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

type UrlRepository interface {
	GetUrl(ctx context.Context, code string) (Shorturl, error)
	CreateUrl(ctx context.Context, params CreateUrlParams) (Shorturl, error)
	ListUrls(ctx context.Context) ([]Shorturl, error)
	DeleteUrl(ctx context.Context, code string) error
	Close()
}

type PostgresRepo struct {
	*Queries
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
		Queries: New(pool),
		pool: pool,
	}, nil
}
