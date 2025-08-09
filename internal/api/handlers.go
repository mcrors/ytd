package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func (s *Server) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	log.Printf("download request: URL=%s, TargetDir=%s, NewName=%s",
		req.URL, req.TargetDir, req.NewName,
	)

	if err := s.dl.Download(r.Context(), req.URL, req.TargetDir, req.NewName); err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("download failed: %v", err))
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Download started"})
}

func (s *Server) GetDirectoriesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("get directories request received")

	entries, err := os.ReadDir(s.baseDir)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	dirs := findDirs(entries)
	resp := DirectoriesResponse{Directories: dirs}
	respondJSON(w, http.StatusOK, resp)
}

func (s *Server) CreateDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req CreateDirectoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	dirPath := filepath.Join(s.baseDir, req.Dir)
	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("created directory: %s", req.Dir)
	respondJSON(w, http.StatusCreated, map[string]string{"message": "Directory created"})
}
