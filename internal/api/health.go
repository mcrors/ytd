package api

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type ReadyConfig struct {
	// BaseDir is where downloads are written; pulled from env with a default by caller.
	BaseDir string
	// YtDlpBin is the executable name or absolute path (default: "yt-dlp").
	YtDlpBin string
	// CommandTimeout bounds external checks like `yt-dlp --version`.
	CommandTimeout time.Duration
}

type readyChecks map[string]string

type healthStatus string

const (
	StatusOK       healthStatus = "ok"
	StatusDegraded healthStatus = "degraded"
)

type status struct {
	Status healthStatus `json:"status"` // "ok" or "degraded"
	Checks readyChecks  `json:"checks"`
}

func ensureWritable(dir string) error {
	if dir == "" {
		return errors.New("empty baseDir")
	}
	// Create the directory tree if missing.
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	// Attempt to create and remove a temp file to prove writability.
	f, err := os.CreateTemp(dir, ".probe-*")
	if err != nil {
		return err
	}
	path := f.Name()
	_ = f.Close()
	_ = os.Remove(path)
	// Resolve symlinks & normalize (defensive; optional).
	_, err = filepath.EvalSymlinks(dir)
	return err
}

func checkYtDlp(parent context.Context, bin string, timeout time.Duration) error {
	// Ensure it's discoverable.
	full, err := exec.LookPath(bin)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, full, "--version")
	// Don’t care about output; we just want a quick, successful return.
	if err := cmd.Run(); err != nil {
		return err
	}
	return ctx.Err() // returns nil unless timed out
}
