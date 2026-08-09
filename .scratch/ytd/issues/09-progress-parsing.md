Status: ready-for-agent
Phase: 2 — Download Flow

Read yt-dlp stdout in real time via `cmd.StdoutPipe()`. Parse progress lines (`[download] 45.3% of 123MiB at 2.34MiB/s ETA 00:45`) and write percentage to the downloads table so the 2s HTMX poll can read it.

## Acceptance criteria

- [ ] `Downloader.Download()` uses `StdoutPipe()` instead of `CombinedOutput()`
- [ ] Progress line parser extracts percentage (integer 0–100)
- [ ] Percentage written to a `progress` column on the downloads row as it changes
- [ ] Non-progress stderr lines captured and stored in `error_message` on failure
- [ ] Unit tests for progress line parser (table-driven, various yt-dlp output formats)
