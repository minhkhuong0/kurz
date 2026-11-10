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
	GetLongUrl(code string) *ShortUrl
	PutShortUrl(url string) *ShortUrl
	SetupDatabase() error
}
