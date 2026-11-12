package main

import (
	"flag"
	"fmt"
	// "forge.nhmk.de/khuong/kurz/internal/handlers"
	// "github.com/gin-gonic/gin"
	"context"
	"os"
	"forge.nhmk.de/khuong/kurz/internal/repository"
	"log"
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

	repo, err := repository.NewRepository(connUrl)
	if err != nil {
		log.Fatalf("Can't start database: %v", err)
	}

	urls, _ := repo.ListUrls(ctx)
	fmt.Print(urls)

	//route()
}
