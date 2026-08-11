package queue

import (
	"context"
	"database/sql"
	"log"
	"sync"

	"github.com/mcrors/ytd/internal/download"
)

const bufferSize = 100

type Downloader interface {
	Download(ctx context.Context, url, targetDir, newName string, format download.Format) error
}

type DownloadJob struct {
	ID        int64
	URL       string
	TargetDir string // absolute path, already validated by the enqueuing layer
	NewName   string
	Format    download.Format
}

type Queue struct {
	dl      Downloader
	db      *sql.DB
	jobs    chan DownloadJob
	workers int
	wg      sync.WaitGroup
}

func New(workers int, db *sql.DB, dl Downloader) *Queue {
	return &Queue{
		dl:      dl,
		db:      db,
		workers: workers,
		jobs:    make(chan DownloadJob, bufferSize),
	}
}

func (q *Queue) Start() {
	for range q.workers {
		q.wg.Add(1)
		go q.worker()
	}
}

// Shutdown stops accepting new jobs and waits for in-flight downloads to finish.
func (q *Queue) Shutdown() {
	close(q.jobs)
	q.wg.Wait()
}

// Enqueue persists the job to SQLite and submits it to the worker pool.
// Returns the assigned download ID.
func (q *Queue) Enqueue(ctx context.Context, job DownloadJob) (int64, error) {
	res, err := q.db.ExecContext(ctx, `
		INSERT INTO downloads (url, target_dir, format, status)
		VALUES (?, ?, ?, 'queued')
	`, job.URL, job.TargetDir, job.Format)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	job.ID = id
	q.jobs <- job
	return id, nil
}

func (q *Queue) worker() {
	defer q.wg.Done()
	for job := range q.jobs {
		q.process(job)
	}
}

func (q *Queue) process(job DownloadJob) {
	if _, err := q.db.Exec(`
		UPDATE downloads SET status='downloading', updated_at=CURRENT_TIMESTAMP WHERE id=?
	`, job.ID); err != nil {
		log.Printf("queue: failed to update status for download %d: %v", job.ID, err)
	}

	err := q.dl.Download(context.Background(), job.URL, job.TargetDir, job.NewName, job.Format)
	if err != nil {
		log.Printf("queue: download %d failed: %v", job.ID, err)
		q.db.Exec(`
			UPDATE downloads SET status='failed', error_message=?, updated_at=CURRENT_TIMESTAMP WHERE id=?
		`, err.Error(), job.ID)
		return
	}

	log.Printf("queue: download %d completed", job.ID)
	q.db.Exec(`
		UPDATE downloads SET status='completed', updated_at=CURRENT_TIMESTAMP WHERE id=?
	`, job.ID)
}
