# ytd — Project Map

Self-hosted YouTube homelab downloader. Go + HTMX + SQLite + yt-dlp. Single user, browser-only, deployed on k3s.

## Decisions so far

- **HTMX-only, no JSON API** — internal/api removed entirely. See ADR-0001.
- **In-memory download queue** — channel-based worker pool, not persisted. Cleared on restart; user resubmits if needed.
- **2s HTMX polling for progress** — no SSE. Simple, debuggable, sufficient for a homelab tool.
- **3 curated download format presets** — Best quality, 1080p, Audio only (MP3). No raw yt-dlp format strings exposed.
- **Title pre-fetched before queuing** — `yt-dlp --get-title` called on submission so queue and history always show a readable title.
- **Cancellation cleans up .part files** — subprocess killed via context; partial files deleted from target dir.
- **Subscriptions: date cutoff + archive file** — `--dateafter <created_at>` prevents historical backfill; `--download-archive` prevents re-download within the window. See ADR-0002.
- **Single global POLL_INTERVAL** — no per-subscription intervals.
- **Poll failures logged, not prominent** — separate poll history view; main subscriptions page stays clean.
- **Library browser is read-only** — no move/rename/delete in v1. Target folder chosen at download time is permanent.

## Phases

1. Foundation cleanup (03–05)
2. Download flow (06–11)
3. UI (12–16)
4. Subscriptions (17–20)
5. Library browser (21–22)
6. Deployment (23)
