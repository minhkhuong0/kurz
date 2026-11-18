package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"forge.nhmk.de/khuong/kurz/internal/db"
	"forge.nhmk.de/khuong/kurz/internal/handlers"
	"github.com/gin-gonic/gin"
)

var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

//go:embed frontend/dist
var frontend embed.FS

func route(h *handlers.Handler) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	router := gin.Default()
	router.GET("/api/url", h.GetUrls)
	router.GET("/api/url/:code", h.GetCode)
	router.POST("/api/url", h.PostUrl)

	frontendFS, err := fs.Sub(frontend, "frontend/dist")
	if err != nil {
		log.Fatalf("Embedded frontend error: %v", err)
	}
	httpFS := http.FS(frontendFS)
	router.GET("/", func(c *gin.Context) {
		c.FileFromFS("/", httpFS)
	})
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.FileFromFS("/favicon.ico", httpFS)
	})
	router.GET("/assets/:name", func(c *gin.Context) {
		c.FileFromFS("assets/"+c.Param("name"), httpFS)
	})

	router.GET(":code", h.Redirect)

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

	route(&handlers.Handler{Repo: repo, DefaultLength: 3, InsertAttempts: 3})
}
