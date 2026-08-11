Status: done
Phase: 2 — Download Flow

In-memory channel-based worker pool for executing yt-dlp downloads concurrently. Bounded by `MAX_CONCURRENT_DOWNLOADS`. Not persisted — queue is cleared on restart.

## Acceptance criteria

- [x] `internal/queue/` package with a `Queue` type
- [x] Accepts `DownloadJob` structs; workers pull from a buffered channel
- [x] Number of workers = `MAX_CONCURRENT_DOWNLOADS` from config
- [x] Each job runs yt-dlp via the `Downloader` interface
- [x] Job state transitions written to SQLite: `queued → downloading → completed | failed | cancelled`
- [x] Queue is started in `main.go` and shut down gracefully on SIGTERM
- [x] Unit tests for worker pool behaviour (mock downloader)
