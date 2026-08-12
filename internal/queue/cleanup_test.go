package queue_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mcrors/ytd/internal/queue"
)

func TestCleanPartFiles(t *testing.T) {
	dir := t.TempDir()

	// Create a mix of .part files and a regular file.
	partFiles := []string{"video.part", "audio.part"}
	otherFiles := []string{"video.mp4"}

	for _, f := range append(partFiles, otherFiles...) {
		if err := os.WriteFile(filepath.Join(dir, f), []byte("data"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	if err := queue.CleanPartFiles(dir); err != nil {
		t.Fatalf("CleanPartFiles: %v", err)
	}

	for _, f := range partFiles {
		if _, err := os.Stat(filepath.Join(dir, f)); !os.IsNotExist(err) {
			t.Errorf(".part file %s should have been deleted", f)
		}
	}
	for _, f := range otherFiles {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Errorf("non-.part file %s should still exist: %v", f, err)
		}
	}
}

func TestCleanPartFiles_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	if err := queue.CleanPartFiles(dir); err != nil {
		t.Fatalf("CleanPartFiles on empty dir: %v", err)
	}
}
