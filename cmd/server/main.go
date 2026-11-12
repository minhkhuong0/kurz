package main

import (
	"flag"
	"fmt"
	// "forge.nhmk.de/khuong/kurz/internal/handlers"
	// "github.com/gin-gonic/gin"
	"context"
	"os"
	"forge.nhmk.de/khuong/kurz/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"errors"
	"github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

// func route() {
// 	router := gin.Default()
// 	router.GET("/api/urls", handlers.GetUrls)
// 	router.GET(":code", handlers.Redirect)
// 	router.POST("/api/urls", handlers.PostUrl)
//
// 	router.Run(":4392")
// }

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

	ctx := context.Background()

	connUrl := os.Getenv("POSTGRESQL_URL")
	pool, err := pgxpool.New(context.Background(), connUrl)
	if err != nil {
		log.Fatalf("Can't connect to database: %v", err)
		os.Exit(1)
	}
	defer pool.Close()

	m, err := migrate.New(
		"file://db/postgres/migrations",
		connUrl,
	)
	if err != nil {
		log.Fatalf("Migration init error: %v", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migration error: %v", err)
	}


	repo := repository.New(pool)
	urls, _ := repo.ListUrls(ctx)
	fmt.Print(urls)

	//route()
}
