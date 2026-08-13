package web

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mcrors/ytd/internal/format"
	"github.com/mcrors/ytd/internal/pathutil"
)

// --- Request/response types ---

type createDirectoryRequest struct {
	Dir string `json:"dir"`
}

type directoriesResponse struct {
	Directories []string `json:"directories"`
}

// --- Readyz types ---

type readyChecks map[string]string

type healthStatus string

const (
	statusOK       healthStatus = "ok"
	statusDegraded healthStatus = "degraded"
)

type readyStatus struct {
	Status healthStatus `json:"status"`
	Checks readyChecks  `json:"checks"`
}

// --- Handlers ---

func (s *server) indexHandler(w http.ResponseWriter, r *http.Request) {
	s.render(w, "index.html", nil)
}

func (s *server) submitHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		respondError(w, http.StatusBadRequest, "invalid form")
		return
	}

	rawURL := r.FormValue("url")
	if rawURL == "" {
		respondError(w, http.StatusBadRequest, "url is required")
		return
	}

	absDir, err := pathutil.SafeJoin(s.baseDir, r.FormValue("targetDir"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid directory: "+err.Error())
		return
	}

	f := format.Format(r.FormValue("format"))
	if f == "" {
		f = format.FormatBest
	}

	title, err := s.queue.GetTitle(r.Context(), rawURL)
	if err != nil {
		log.Printf("web: GetTitle: %v", err)
		title = rawURL
	}

	id, err := s.queue.Enqueue(r.Context(), rawURL, absDir, r.FormValue("newName"), f, title)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "enqueue failed: "+err.Error())
		return
	}

	d, err := queryDownload(s.db, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.renderFragment(w, "download-row.html", d)
}

func (s *server) statusHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid download id")
		return
	}
	d, err := queryDownload(s.db, id)
	if errors.Is(err, sql.ErrNoRows) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.renderFragment(w, "download-row.html", d)
}

func (s *server) historyHandler(w http.ResponseWriter, r *http.Request) {
	downloads, err := queryHistory(s.db)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.renderFragment(w, "download-history.html", downloads)
}

func (s *server) getDirectoriesHandler(w http.ResponseWriter, r *http.Request) {
	entries, err := os.ReadDir(s.baseDir)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, directoriesResponse{Directories: findDirs(entries)})
}

func (s *server) createDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req createDirectoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	target, err := pathutil.SafeJoin(s.baseDir, req.Dir)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := os.MkdirAll(target, 0o755); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("created directory: %s", req.Dir)
	respondJSON(w, http.StatusCreated, map[string]string{"message": "Directory created"})
}

func (s *server) cancelHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid download id")
		return
	}
	s.queue.Cancel(id)
	w.WriteHeader(http.StatusNoContent)
}

func (s *server) healthzHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *server) readyzHandler(w http.ResponseWriter, r *http.Request) {
	checks := readyChecks{"baseDir": "ok", "db": "ok", "yt-dlp": "ok"}
	var readyErr error

	if err := ensureWritable(s.baseDir); err != nil {
		checks["baseDir"] = err.Error()
		readyErr = errors.New("baseDir not writable")
	}
	if err := s.db.PingContext(r.Context()); err != nil {
		checks["db"] = err.Error()
		readyErr = errors.New("db unavailable")
	}
	if err := checkYtDlp(r.Context(), "yt-dlp", 2*time.Second); err != nil {
		checks["yt-dlp"] = err.Error()
		readyErr = errors.New("yt-dlp unavailable")
	}

	st, code := statusOK, http.StatusOK
	if readyErr != nil {
		st, code = statusDegraded, http.StatusServiceUnavailable
	}
	respondJSON(w, code, readyStatus{Status: st, Checks: checks})
}

// --- Health helpers ---

func ensureWritable(dir string) error {
	if dir == "" {
		return errors.New("empty baseDir")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	f, err := os.CreateTemp(dir, ".probe-*")
	if err != nil {
		return err
	}
	path := f.Name()
	_ = f.Close()
	_ = os.Remove(path)
	_, err = filepath.EvalSymlinks(dir)
	return err
}

func checkYtDlp(parent context.Context, bin string, timeout time.Duration) error {
	full, err := exec.LookPath(bin)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()
	if err := exec.CommandContext(ctx, full, "--version").Run(); err != nil {
		return err
	}
	return ctx.Err()
}

func findDirs(entries []os.DirEntry) []string {
	var results []string
	for _, e := range entries {
		if e.IsDir() {
			results = append(results, e.Name())
		}
	}
	return results
}
