package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func (s *Server) downloadHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	log.Printf("download request: URL=%s, TargetDir=%s, NewName=%s",
		req.URL, req.TargetDir, req.NewName,
	)

	rel, err := normalizeTwoLevel(req.TargetDir)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	target := filepath.Join(s.baseDir, rel)

	if err := s.dl.Download(r.Context(), req.URL, target, req.NewName); err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("download failed: %v", err))
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Download started"})
}

func (s *Server) getDirectoriesHandler(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) createDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req CreateDirectoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	rel, err := normalizeTwoLevel(req.Dir)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	dirPath := filepath.Join(s.baseDir, rel)

	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("created directory: %s", req.Dir)
	respondJSON(w, http.StatusCreated, map[string]string{"message": "Directory created"})
}

func (s *Server) healthzHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) readyzHandlers(w http.ResponseWriter, r *http.Request) {
	checks := readyChecks{
		"baseDir": "ok",
		"yt-dlp":  "ok",
	}

	var readyErr error

	if err := ensureWritable(s.baseDir); err != nil {
		checks["baseDir"] = err.Error()
		readyErr = errors.New("baseDir not writable")
	}

	bin := "yt-dlp"
	timeout := 2 * time.Second
	if err := checkYtDlp(r.Context(), bin, timeout); err != nil {
		checks["yt-dlp"] = err.Error()
		readyErr = errors.New("yt-dlp unavailable")
	}

	resp := status{
		Status: map[bool]healthStatus{true: StatusOK, false: StatusDegraded}[readyErr == nil],
		Checks: checks,
	}

	statusCode := http.StatusOK
	if readyErr != nil {
		statusCode = http.StatusServiceUnavailable // 503 so k8s/infra can gate traffic
	}
	respondJSON(w, statusCode, resp)
}
