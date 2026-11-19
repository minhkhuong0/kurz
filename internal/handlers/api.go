package handlers

import (
	"context"
	"fmt"
	"forge.nhmk.de/khuong/kurz/internal/db"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Handler struct {
	Repo           db.UrlRepository
	DefaultLength  int
	InsertAttempts int
}

type postUrlRequest struct {
	Url string `json:"url" binding:"required"`
	Ttl string `json:"ttl"`
}

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"0123456789"

func generateCode(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
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

func (h *Handler) GetCode(c *gin.Context) {
	code := c.Param("code")
	url, err := h.retrieveCode(code)
	if err != nil {
		c.String(http.StatusNotFound, "404 code not found\n")
		return
	}
	c.IndentedJSON(http.StatusOK, url)
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
	for i := 0; i < h.InsertAttempts; i++ {
		arg := db.CreateUrlParams{
			Code:      generateCode(h.DefaultLength),
			Url:       rq.Url,
			CreatedAt: time.Now(),
			ExpiresAt: time.Now(),
		}
		err := h.Repo.CreateUrl(ctx, arg)
		if err == nil {
			c.IndentedJSON(http.StatusCreated, arg)
			return
		} else if i == h.InsertAttempts-1 {
			c.IndentedJSON(http.StatusServiceUnavailable, gin.H{
				"error": "Can't find free short code in database. Please try again.",
			})
			log.Printf("Error inserting to database: %v", err)
			return
		}
	}
}

func (h *Handler) retrieveCode(code string) (db.Shorturl, error) {
	ctx := context.Background()
	url, err := h.Repo.GetUrl(ctx, code)
	if err != nil {
		return db.Shorturl{}, fmt.Errorf("Code not in database: %w", err)
	}
	return url, nil
}

func (h *Handler) Redirect(c *gin.Context) {
	code := c.Param("code")
	url, err := h.retrieveCode(code)
	if err != nil {
		c.String(http.StatusNotFound, "404 page not found\n")
		return
	}
	c.Redirect(http.StatusMovedPermanently, url.Url)
}
