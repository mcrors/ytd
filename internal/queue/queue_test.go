package queue_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/mcrors/ytd/internal/db"
	"github.com/mcrors/ytd/internal/queue"
)

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	// Single connection prevents each goroutine from opening a fresh :memory: db.
	database, err := db.Open(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	database.SetMaxOpenConns(1)
	if err := db.Migrate(database); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	t.Cleanup(func() { database.Close() })
	return database
}

// funcDownloader lets tests supply a Download implementation inline.
type funcDownloader struct {
	fn func(ctx context.Context, url, dir, name string) error
}

func (f *funcDownloader) Download(ctx context.Context, url, dir, name string) error {
	return f.fn(ctx, url, dir, name)
}

func okDownloader() *funcDownloader {
	return &funcDownloader{fn: func(_ context.Context, _, _, _ string) error { return nil }}
}

func errDownloader(msg string) *funcDownloader {
	return &funcDownloader{fn: func(_ context.Context, _, _, _ string) error {
		return errors.New(msg)
	}}
}

func jobStatus(t *testing.T, database *sql.DB, id int64) (status, errMsg string) {
	t.Helper()
	if err := database.QueryRow(
		`SELECT status, error_message FROM downloads WHERE id=?`, id,
	).Scan(&status, &errMsg); err != nil {
		t.Fatalf("querying status for download %d: %v", id, err)
	}
	return
}

func TestQueue_CompletedJob(t *testing.T) {
	database := newTestDB(t)
	q := queue.New(1, database, okDownloader())
	q.Start()

	id, err := q.Enqueue(context.Background(), queue.DownloadJob{
		URL:       "https://example.com/video",
		TargetDir: "/tmp/media",
	})
	if err != nil {
		t.Fatalf("Enqueue: %v", err)
	}

	q.Shutdown()

	status, _ := jobStatus(t, database, id)
	if status != "completed" {
		t.Errorf("status = %q, want completed", status)
	}
}

func TestQueue_FailedJob(t *testing.T) {
	database := newTestDB(t)
	q := queue.New(1, database, errDownloader("yt-dlp exploded"))
	q.Start()

	id, err := q.Enqueue(context.Background(), queue.DownloadJob{
		URL:       "https://example.com/video",
		TargetDir: "/tmp/media",
	})
	if err != nil {
		t.Fatalf("Enqueue: %v", err)
	}

	q.Shutdown()

	status, errMsg := jobStatus(t, database, id)
	if status != "failed" {
		t.Errorf("status = %q, want failed", status)
	}
	if errMsg == "" {
		t.Error("want non-empty error_message")
	}
}

func TestQueue_BoundedConcurrency(t *testing.T) {
	const workers = 2
	const jobs = 6

	database := newTestDB(t)

	var mu sync.Mutex
	current, maxSeen := 0, 0

	dl := &funcDownloader{fn: func(_ context.Context, _, _, _ string) error {
		mu.Lock()
		current++
		if current > maxSeen {
			maxSeen = current
		}
		mu.Unlock()

		time.Sleep(20 * time.Millisecond)

		mu.Lock()
		current--
		mu.Unlock()
		return nil
	}}

	q := queue.New(workers, database, dl)
	q.Start()

	for i := range jobs {
		if _, err := q.Enqueue(context.Background(), queue.DownloadJob{
			URL:       fmt.Sprintf("https://example.com/video/%d", i),
			TargetDir: "/tmp/media",
		}); err != nil {
			t.Fatalf("Enqueue %d: %v", i, err)
		}
	}

	q.Shutdown()

	if maxSeen > workers {
		t.Errorf("max concurrent = %d, want <= %d", maxSeen, workers)
	}
	if maxSeen == 0 {
		t.Error("no jobs ran concurrently — workers may not have started")
	}
}
