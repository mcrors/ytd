package api

import (
	"context"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type DownloadService interface {
	Download(ctx context.Context, url, targetDir, newName string) error
}

type Server struct {
	dl      DownloadService
	baseDir string
}

func NewServer(dl DownloadService) http.Handler {
	base := getenv("YTD_BASE_DIR", "./data/media/youtube")

	s := &Server{dl: dl, baseDir: base}

	r := mux.NewRouter()

	r.HandleFunc("/api/download", s.DownloadHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/directories", s.GetDirectoriesHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/directory", s.CreateDirectoryHandler).Methods(http.MethodPost)

	return r
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
