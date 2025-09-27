package downloader

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Commander func(ctx context.Context, name string, args ...string) *exec.Cmd
type LookPathFunc func(file string) (string, error)

type YouTube struct {
	// which binary to call; defaults to "yt-dlp"
	Bin      string
	Command  Commander
	LookPath LookPathFunc
}

func NewYouTube() *YouTube {
	return &YouTube{
		Bin:      "yt-dlp",
		Command:  exec.CommandContext,
		LookPath: exec.LookPath,
	}
}

// Download downloads a video from the given YouTube URL.
//
// Parameters:
//   - ctx: context used to control cancellation or timeouts for the command execution.
//   - url: the YouTube video (or playlist) URL to download.
//   - targetDir: the directory where the downloaded file should be saved. It will
//     be created if it does not exist.
//   - newName: optional new base name for the output file. If empty, the video title
//     is used instead.
//
// Returns:
//   - error: non-nil if the binary is not found in PATH, the target directory
//     cannot be created, or the download command fails.
func (y *YouTube) Download(ctx context.Context, url, targetDir, newName string) error {
	if _, err := y.LookPath(y.Bin); err != nil {
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

	// TODO: should this run on parallel
	cmd := y.Command(ctx, y.Bin, "-o", filepath.Join(targetDir, outTpl), url)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s failed: %w\n%s", y.Bin, err, string(out))
	}
	return nil
}

// GetChannel retrieves the channel name for a given YouTube URL.
//
// Parameters:
//   - ctx: context used to control cancellation or timeouts for the command execution.
//   - url: the YouTube video or channel URL to inspect.
//
// Returns:
//   - string: the resolved channel name, with leading and trailing whitespace removed.
//   - error: non-nil if the binary is not found in PATH, or if the command execution fails.
func (y *YouTube) GetChannel(ctx context.Context, url string) (string, error) {
	if _, err := y.LookPath(y.Bin); err != nil {
		return "", fmt.Errorf("%s not found in PATH: %w", y.Bin, err)
	}
	cmd := y.Command(ctx, y.Bin, "--no-warning", "--print", "%(channel)s", url)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s failed: %w\n%s", y.Bin, err, string(out))
	}
	result := strings.TrimSpace(string(out))
	return result, nil
}
