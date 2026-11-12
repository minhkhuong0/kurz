package tools

import (
	"time"
)

type ShortUrl struct {
	Code      string    `json:"code"`
	Url       string    `json:"url"`
	ExpiresAt time.Time `json:"expires_at"`
}

type DataBaseInterface interface {
	GetUrl(code string) *ShortUrl
	PutUrl(url string) *ShortUrl
	SetupDatabase() error
}

func SetupDatabase() error {
	db, err := sql.Open(v5,)
}
