package downloader

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type YouTube struct {
	// which binary to call; defaults to "yt-dlp"
	Bin string
}

func NewYouTube() *YouTube {
	return &YouTube{Bin: "yt-dlp"}
}

func (y *YouTube) Download(ctx context.Context, url, targetDir, newName string) error {
	if _, err := exec.LookPath(y.Bin); err != nil {
		return fmt.Errorf("%s not found in PATH: %w", y.Bin, err)
	}
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	outTpl := "%(title)s.%(ext)s"
	if newName != "" {
		newName = filepath.Base(newName)
		outTpl = newName + ".%(ext)s"
	}

	cmd := exec.CommandContext(ctx, y.Bin, "-o", filepath.Join(targetDir, outTpl), url)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s failed: %w\n%s", y.Bin, err, string(out))
	}
	return nil
}
