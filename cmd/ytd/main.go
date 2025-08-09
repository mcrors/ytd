package main

import (
	"log"
	"net/http"

	"github.com/mcrors/ytd/internal/api"
	"github.com/mcrors/ytd/internal/downloader"
)

func main() {
	yt := downloader.NewYouTube()
	server := api.NewServer(yt)

	log.Println("server running on port 8080 ...")
	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatal(err)
	}
}
