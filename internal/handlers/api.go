package handlers

import (
	"forge.nhmk.de/khuong/kurz/internal/tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var shortUrls = []tools.ShortUrl{
	{Code: "abc", Url: "https://nhmk.de", ExpiresAt: time.Now()},
	{Code: "aaa", Url: "https://github.com", ExpiresAt: time.Now()},
	{Code: "zxc", Url: "https://example.org", ExpiresAt: time.Now()},
}

func GetUrls(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, shortUrls)
}
