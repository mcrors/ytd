package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewServer() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/api/download", DownloadHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/directories", GetDirectoriesHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/directory", CreateDirectoryHandler).Methods(http.MethodPost)

	return r
}
