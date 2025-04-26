package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mcrors/ytd/downloader"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("download request received: URL=%s, TargetDir=%s, NewName=%s\n",
		req.URL, req.TargetDir, req.NewName,
	)
	downloader.DownloadVideo(req.URL, req.TargetDir, req.NewName)
	respondJson(w, http.StatusOK, map[string]string{"message": "Download request received"})
}

func GetDirectoriesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("get directories request received")
	data := []string{"tech", "music", "history"}
	resp := DirectoriesResponse{Directories: data}
	respondJson(w, http.StatusOK, resp)
}

func CreateDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req CreateDirectoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("creating directory: %s\n", req.Dir)
	respondJson(w, http.StatusCreated, map[string]string{"message": "directory created"})
}
