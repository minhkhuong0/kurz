package main

import (
	"forge.nhmk.de/khuong/kurz/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/api/urls", handlers.GetUrls)

	router.Run("localhost:8080")
}
