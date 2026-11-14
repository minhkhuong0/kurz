package db

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"forge.nhmk.de/khuong/kurz/internal/repository"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
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

//go:embed postgres/migrations/*.sql
var fs embed.FS

func NewRepository(connUrl string) (UrlRepository, error) {
	pool, err := pgxpool.New(context.Background(), connUrl)
	if err != nil {
		return nil, fmt.Errorf("Can't establish a connection pool to the database %w", err)
	}

	d, err := iofs.New(fs, "postgres/migrations")
	if err != nil {
		return nil, fmt.Errorf("Migration source error: %w", err)
	}
	m, err := migrate.NewWithSourceInstance(
		"iofs", d,
		connUrl,
	)
	if err != nil {
		return nil, fmt.Errorf("Migration init error: %w", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("Migration error: %w", err)
	}

	return &PostgresRepo{
		Queries: repository.New(pool),
		pool:    pool,
	}, nil
}
