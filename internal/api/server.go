package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mcrors/ytd/internal/download"
)

type Downloader interface {
	Download(context.Context, download.DownloadCommand) (*download.DownloadResult, error)
}

type Server struct {
	dl      Downloader
	baseDir string
}

func NewServer(dl Downloader, baseDir string) http.Handler {
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
