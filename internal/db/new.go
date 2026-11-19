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
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type UrlRepository interface {
	GetUrl(ctx context.Context, code string) (Shorturl, error)
	CreateUrl(ctx context.Context, params CreateUrlParams) error
	ListUrls(ctx context.Context) ([]Shorturl, error)
	DeleteUrl(ctx context.Context, code string) error
	Close()
}

type PostgresRepo struct {
	*repository.Queries
	pool *pgxpool.Pool
}

type CreateUrlParams struct {
	Code      string    `json:"code"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (repo *PostgresRepo) Close() {
	repo.pool.Close()
}

type Shorturl struct {
	Code      string    `json:"code"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func toTimestamp(t time.Time) pgtype.Timestamp {
	if t.IsZero() {
		return pgtype.Timestamp{Valid: false}
	}
	return pgtype.Timestamp{Time: t, Valid: true}
}

func toTime(t pgtype.Timestamp) time.Time {
	if t.Valid {
		return t.Time
	}
	return time.Time{}
}

func toDbShorturl(pgUrl repository.Shorturl) Shorturl {
	return Shorturl{
		Code:      pgUrl.Code,
		Url:       pgUrl.Url,
		CreatedAt: toTime(pgUrl.CreatedAt),
		ExpiresAt: toTime(pgUrl.ExpiresAt),
	}
}

func (repo *PostgresRepo) GetUrl(ctx context.Context, code string) (Shorturl, error) {
	pgUrl, err := repo.Queries.GetUrl(ctx, code)
	if err != nil {
		return Shorturl{}, err
	}
	return toDbShorturl(pgUrl), nil
}

func (repo *PostgresRepo) CreateUrl(ctx context.Context, arg CreateUrlParams) error {
	err := repo.Queries.CreateUrl(ctx, repository.CreateUrlParams{
		Code:      arg.Code,
		Url:       arg.Url,
		CreatedAt: toTimestamp(arg.CreatedAt),
		ExpiresAt: toTimestamp(arg.ExpiresAt),
	})
	if err != nil {
		return err
	}
	return nil
}

func (repo *PostgresRepo) ListUrls(ctx context.Context) ([]Shorturl, error) {
	pgUrls, err := repo.Queries.ListUrls(ctx)
	if err != nil {
		return []Shorturl{}, err
	}
	urls := make([]Shorturl, 0, len(pgUrls))
	for _, pgUrl := range pgUrls {
		urls = append(urls, toDbShorturl(pgUrl))
	}
	return urls, nil
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
