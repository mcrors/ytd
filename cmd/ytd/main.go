package main

import (
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/mcrors/ytd/internal/api"
	"github.com/mcrors/ytd/internal/config"
	"github.com/mcrors/ytd/internal/download"
	"github.com/mcrors/ytd/internal/downloader"
	"github.com/mcrors/ytd/internal/middleware"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	yt := downloader.NewYouTube("yt-dlp", exec.CommandContext, exec.LookPath)
	ds := download.NewDownloadService(cfg.MediaDir, yt)
	server := api.NewServer(ds, cfg.MediaDir)

	server = middleware.Logging(server)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           server,
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("server running on port %s ...", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
