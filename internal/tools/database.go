package tools

import (
	"context"
)

type UrlReposityory interface {
	GetUrl(ctx context.Context, code string) (string, error)
	CreateUrl(ctx context.Context, url string) (string, error)
	DeleteUrl(ctx context.Context, code string) (string, error)
}

