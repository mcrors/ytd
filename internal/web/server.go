package web

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/mcrors/ytd/internal/download"
)

type Downloader interface {
	Download(context.Context, download.DownloadCommand) (*download.DownloadResult, error)
}

type server struct {
	dl      Downloader
	baseDir string
	db      *sql.DB
}

func RegisterRoutes(mux *http.ServeMux, dl Downloader, baseDir string, db *sql.DB) {
	s := &server{dl: dl, baseDir: baseDir, db: db}

	mux.HandleFunc("GET /healthz", s.healthzHandler)
	mux.HandleFunc("GET /readyz", s.readyzHandler)
	mux.HandleFunc("POST /api/download", s.downloadHandler)
	mux.HandleFunc("GET /api/directories", s.getDirectoriesHandler)
	mux.HandleFunc("POST /api/directory", s.createDirectoryHandler)
}
