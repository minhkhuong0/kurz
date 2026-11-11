package handlers

import (
	"errors"
	"forge.nhmk.de/khuong/kurz/internal/tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type CreateUrlRequest struct {
	Url string `json:"url"`
}

var shortUrls = []tools.ShortUrl{
	{Code: "abc", Url: "https://nhmk.de", ExpiresAt: time.Now()},
	{Code: "aaa", Url: "https://github.com", ExpiresAt: time.Now()},
	{Code: "zxc", Url: "https://example.org", ExpiresAt: time.Now()},
}

func GetUrls(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, shortUrls)
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

func GetUrl(code string) (string, error) {
	for _, v := range shortUrls {
		if v.Code == code {
			return v.Url, nil
		}
	}
	return "", errors.New("Code not found")
}

func Redirect(c *gin.Context) {
	code := c.Param("code")
	url, err := GetUrl(code)
	if err != nil {
		c.String(http.StatusNotFound, "404 page not found\n")
	}
	c.Redirect(http.StatusMovedPermanently, url)
}
