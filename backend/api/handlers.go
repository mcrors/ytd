package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mcrors/ytd/downloader"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("download request received: URL=%s, TargetDir=%s, NewName=%s\n",
		req.URL, req.TargetDir, req.NewName,
	)
	if err := downloader.DownloadVideo(req.URL, req.TargetDir, req.NewName); err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Download failed: %v", err))
		return
	}
	respondJson(w, http.StatusOK, map[string]string{"message": "Download request received"})
}

func GetDirectoriesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("get directories request received")
	baseDir := "./data/media/youtube" // TODO: make configurable
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	dirs := findDirs(entries)
	resp := DirectoriesResponse{Directories: dirs}
	respondJson(w, http.StatusOK, resp)
}

func CreateDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req CreateDirectoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	baseDir := "./data/media/youtube" // TODO: make configurable
	dirPath := filepath.Join(baseDir, req.Dir)

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("successfully created directory: %s\n", req.Dir)
	respondJson(w, http.StatusCreated, map[string]string{"message": "directory created"})
}
