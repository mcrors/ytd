package downloader

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const baseDir = "./data/media/youtube" // TODO: Make configurable later

// DownloadVideo downloads a video from the given URL to the specified target directory.
// If newName is provided, the video will be saved with that name.
func DownloadVideo(url, targetDir, newName string) error {
	fullTargetDir := filepath.Join(baseDir, targetDir)

	// Ensure the directory exists
	if err := os.MkdirAll(fullTargetDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	outputPath := filepath.Join(fullTargetDir, "%(title)s.%(ext)s")
	if newName != "" {
		outputPath = filepath.Join(fullTargetDir, newName+".%(ext)s")
	}

	cmd := exec.Command("yt-dlp", "-o", outputPath, url)

	// Optionally, you can capture the output for logging
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to download video: %w", err)
	}

	return nil
}
