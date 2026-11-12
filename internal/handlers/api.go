package handlers

import (
	//"errors"
	"forge.nhmk.de/khuong/kurz/internal/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
	"context"
)

type CreateUrlRequest struct {
	Url string `json:"url"`
}

type Handler struct {
	Repo db.UrlRepository
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
	c.JSON(http.StatusOK, urls)
}

func PostUrl(c *gin.Context) {
	var rq CreateUrlRequest
	if err := c.BindJSON(&rq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON: " + err.Error(),
		})
		return
	}

	if rq.Url == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: Field 'url' is required",
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, rq)
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
