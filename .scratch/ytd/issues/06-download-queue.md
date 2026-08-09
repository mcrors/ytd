Status: ready-for-agent
Phase: 2 — Download Flow

In-memory channel-based worker pool for executing yt-dlp downloads concurrently. Bounded by `MAX_CONCURRENT_DOWNLOADS`. Not persisted — queue is cleared on restart.

## Acceptance criteria

- [ ] `internal/queue/` package with a `Queue` type
- [ ] Accepts `DownloadJob` structs; workers pull from a buffered channel
- [ ] Number of workers = `MAX_CONCURRENT_DOWNLOADS` from config
- [ ] Each job runs yt-dlp via the `Downloader` interface
- [ ] Job state transitions written to SQLite: `queued → downloading → completed | failed | cancelled`
- [ ] Queue is started in `main.go` and shut down gracefully on SIGTERM
- [ ] Unit tests for worker pool behaviour (mock downloader)
