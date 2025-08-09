package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mcrors/ytd/internal/api"
	"github.com/mcrors/ytd/internal/downloader"
	"github.com/mcrors/ytd/internal/middleware"
)

func main() {
	baseDir := os.Getenv("YTD_BASE_DIR")
	if baseDir == "" {
		baseDir = "./data/media/youtube"
	}
	yt := downloader.NewYouTube()
	server := api.NewServer(yt, baseDir)

	server = middleware.Logging(server)

	log.Println("server running on port 8080 ...")
	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatal(err)
	}
}
