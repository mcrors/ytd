package downloader

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/mcrors/ytd/internal/download"
)

type Commander func(ctx context.Context, name string, args ...string) *exec.Cmd
type LookPathFunc func(file string) (string, error)

type youTube struct {
	// which binary to call; defaults to "yt-dlp"
	bin          string
	cmd          Commander
	lookPathFunc LookPathFunc
}

func NewYouTube(bin string, cmd Commander, LookPathFunc LookPathFunc) *youTube {
	return &youTube{
		bin:          bin,
		cmd:          cmd,
		lookPathFunc: LookPathFunc,
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
func (y *youTube) Download(ctx context.Context, url, targetDir, newName string, format download.Format, onProgress func(int)) error {
	if _, err := y.lookPathFunc(y.bin); err != nil {
		return fmt.Errorf("%s not found in PATH: %w", y.bin, err)
	}

	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	formatArgs, err := download.FormatArgs(format)
	if err != nil {
		return err
	}

	outTpl := "%(title)s.%(ext)s"
	if newName != "" {
		newName = filepath.Base(newName)
		outTpl = newName + ".%(ext)s"
	}

	args := append(formatArgs, "-o", filepath.Join(targetDir, outTpl), url)
	cmd := y.cmd(ctx, y.bin, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting %s: %w", y.bin, err)
	}

	var stderrBuf bytes.Buffer
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		io.Copy(&stderrBuf, stderr)
	}()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		if pct, ok := ParseProgress(scanner.Text()); ok && onProgress != nil {
			onProgress(pct)
		}
	}

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("%s failed: %w\n%s", y.bin, err, stderrBuf.String())
	}
	return nil
}

// GetTitle retrieves the video title for a given YouTube URL.
func (y *youTube) GetTitle(ctx context.Context, url string) (string, error) {
	if _, err := y.lookPathFunc(y.bin); err != nil {
		return "", fmt.Errorf("%s not found in PATH: %w", y.bin, err)
	}
	cmd := y.cmd(ctx, y.bin, "--no-warning", "--print", "%(title)s", url)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s failed: %w\n%s", y.bin, err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
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
func (y *youTube) GetChannel(ctx context.Context, url string) (string, error) {
	if _, err := y.lookPathFunc(y.bin); err != nil {
		return "", fmt.Errorf("%s not found in PATH: %w", y.bin, err)
	}
	cmd := y.cmd(ctx, y.bin, "--no-warning", "--print", "%(channel)s", url)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s failed: %w\n%s", y.bin, err, string(out))
	}
	result := strings.TrimSpace(string(out))
	return result, nil
}
