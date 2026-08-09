Status: ready-for-agent
Phase: 2 — Download Flow

Before a download is queued, call `yt-dlp --get-title <url>` to fetch the video title and store it immediately in the downloads table. The queue and history always show a human-readable title, never a raw URL.

## Acceptance criteria

- [ ] `Downloader.GetTitle(ctx, url) (string, error)` method (similar shape to existing `GetChannel`)
- [ ] Called in the submit handler before inserting the downloads row
- [ ] Title stored in `downloads.title` at insert time
- [ ] If `--get-title` fails, submission returns an error (don't queue a nameless download)
- [ ] Unit tests with mock downloader
