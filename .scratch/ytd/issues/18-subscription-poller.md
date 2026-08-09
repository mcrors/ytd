Status: ready-for-agent
Phase: 4 — Subscriptions

Background goroutine that ticks on `POLL_INTERVAL`, fetches all active subscriptions from SQLite, and runs yt-dlp for each. Uses `--dateafter <created_at>` to skip historical videos and `--download-archive <path>` to prevent re-downloads. New videos are pushed into the download queue.

## Acceptance criteria

- [ ] `internal/poller/` package with `Poller` type
- [ ] Starts in `main.go`; shuts down gracefully on SIGTERM via context cancellation
- [ ] Ticks every `POLL_INTERVAL`; queries subscriptions where `status = 'active'`
- [ ] Per subscription: calls yt-dlp with `--dateafter`, `--download-archive`, and the subscription's target dir
- [ ] Each new video found enqueues a Download job (via the queue from ticket 06)
- [ ] Updates `last_checked_at` on every poll attempt (success or failure)
- [ ] Updates `last_downloaded_at` when at least one new video is queued
- [ ] On failure: inserts a row into `poll_failures`; does not crash the poller
- [ ] Unit tests with mock downloader and in-memory SQLite
