package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type DownloadService interface {
	Download(ctx context.Context, url, targetDir, newName string) error
}

type Server struct {
	dl      DownloadService
	baseDir string
}

func NewServer(dl DownloadService, baseDir string) http.Handler {
	s := &Server{
		dl:      dl,
		baseDir: baseDir,
	}

	r := mux.NewRouter()

	r.HandleFunc("/healthz", s.healthzHandler).Methods(http.MethodGet)
	r.HandleFunc("/readyz", s.readyzHandlers).Methods(http.MethodGet)

	r.HandleFunc("/api/download", s.downloadHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/directories", s.getDirectoriesHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/directory", s.createDirectoryHandler).Methods(http.MethodPost)

	return r
}
