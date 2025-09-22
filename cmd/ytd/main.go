package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mcrors/ytd/internal/api"
	"github.com/mcrors/ytd/internal/download"
	"github.com/mcrors/ytd/internal/downloader"
	"github.com/mcrors/ytd/internal/middleware"
)

func main() {
	baseDir := os.Getenv("YTD_BASE_DIR")
	if baseDir == "" {
		baseDir = "./data/media/youtube"
	}
	yt := downloader.NewYouTube()
	ds := download.NewDownloadService(baseDir, yt)
	server := api.NewServer(ds, baseDir)

	server = middleware.Logging(server)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           server,
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Println("server running on port 8080 ...")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
