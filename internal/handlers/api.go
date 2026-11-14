package handlers

import (
	//"errors"
	"context"
	"forge.nhmk.de/khuong/kurz/internal/db"
	"forge.nhmk.de/khuong/kurz/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	Repo db.UrlRepository
}

type postUrlRequest struct {
	Url string `json:"url" binding:"required"`
	Ttl string `json:"ttl"`
}

func ToTimestamp(t time.Time) pgtype.Timestamp {
	if t.IsZero() {
		return pgtype.Timestamp{Valid: false}
	}
	return pgtype.Timestamp{Time: t, Valid: true}
}

func (h *Handler) GetUrls(c *gin.Context) {
	ctx := context.Background()
	urls, err := h.Repo.ListUrls(ctx)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal database error",
		})
		log.Printf("Error getting urls from database: %v", err)
		return
	}
	c.IndentedJSON(http.StatusOK, urls)
}

func (h *Handler) PostUrl(c *gin.Context) {
	var rq postUrlRequest
	if err := c.BindJSON(&rq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON: " + err.Error(),
		})
		return
	}
	ctx := context.Background()
	shortUrl, err := h.Repo.CreateUrl(ctx, repository.CreateUrlParams{
		Code:      "aaa",
		Url:       rq.Url,
		CreatedAt: ToTimestamp(time.Now()),
		ExpiresAt: ToTimestamp(time.Now()),
	})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal database error",
		})
		log.Printf("Error getting created url from database: %v", err)
		return
	}

	c.IndentedJSON(http.StatusCreated, shortUrl)
}

// func GetUrl(code string) (string, error) {
// 	for _, v := range shortUrls {
// 		if v.Code == code {
// 			return v.Url, nil
// 		}
// 	}
// 	return "", errors.New("Code not found")
// }
//
// func Redirect(c *gin.Context) {
// 	code := c.Param("code")
// 	url, err := GetUrl(code)
// 	if err != nil {
// 		c.String(http.StatusNotFound, "404 page not found\n")
// 		return
// 	}
// 	c.Redirect(http.StatusMovedPermanently, url)
// }
