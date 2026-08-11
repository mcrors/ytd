package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/mcrors/ytd/internal/config"
	"github.com/mcrors/ytd/internal/db"
	"github.com/mcrors/ytd/internal/download"
	"github.com/mcrors/ytd/internal/downloader"
	"github.com/mcrors/ytd/internal/middleware"
	"github.com/mcrors/ytd/internal/queue"
	"github.com/mcrors/ytd/internal/web"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	database, err := db.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer database.Close()

	if err := db.Migrate(database); err != nil {
		log.Fatalf("db migrate: %v", err)
	}

	yt := downloader.NewYouTube("yt-dlp", exec.CommandContext, exec.LookPath)
	ds := download.NewDownloadService(cfg.MediaDir, yt)

	q := queue.New(cfg.MaxConcurrentDL, database, yt)
	q.Start()

	mux := http.NewServeMux()
	if err := web.RegisterRoutes(mux, ds, cfg.MediaDir, database, cfg.Dev); err != nil {
		log.Fatalf("web: %v", err)
	}

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           middleware.Logging(mux),
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("server running on port %s ...", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	<-quit

	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown: %v", err)
	}
	q.Shutdown()
}
