package queue

import (
	"log"
	"os"
	"path/filepath"
)

// CleanPartFiles removes any *.part files left in dir by an interrupted yt-dlp run.
func CleanPartFiles(dir string) error {
	matches, err := filepath.Glob(filepath.Join(dir, "*.part"))
	if err != nil {
		return err
	}
	for _, f := range matches {
		if err := os.Remove(f); err != nil {
			log.Printf("queue: removing part file %s: %v", f, err)
		}
	}
	return nil
}
