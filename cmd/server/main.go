package main

import (
	"flag"
	"fmt"
	"forge.nhmk.de/khuong/kurz/internal/handlers"
	"github.com/gin-gonic/gin"
	"os"
	"forge.nhmk.de/khuong/kurz/internal/db"
	"log"
)

var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

func route(h *handlers.Handler) {
	router := gin.Default()
	router.GET("/api/urls", h.GetUrls)
	router.POST("/api/urls", handlers.PostUrl)

	// router.GET(":code", h.Redirect)

	router.Run(":4392")
}

func main() {
	const (
		usageVersion = "Print version and exit"
	)
	versionFlag := flag.Bool("version", false, usageVersion)
	versionFlag = flag.Bool("v", false, usageVersion+" (Shorhand)")
	flag.Parse()


	if *versionFlag {
		fmt.Printf("Kurz URL Shortener\nVersion: %s\nCommit: %s\nBuilt :%s\n",
			Version, GitCommit, BuildDate)
		return
	}

	connUrl := os.Getenv("POSTGRESQL_URL")

	repo, err := db.NewRepository(connUrl)
	if err != nil {
		log.Fatalf("Can't start database: %v", err)
	}
	defer repo.Close()

	route(&handlers.Handler{ Repo: repo })
}
