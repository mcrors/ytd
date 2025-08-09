# TODO

## Phase 1 — Cleanup & Safety
- [ ] Add `YTD_BASE_DIR` environment variable (default to `./data/media/youtube`).
- [ ] Implement `safeJoin` or `normalizeTwoLevel` to:
  - Prevent path traversal (`../` etc.).
  - Enforce exactly **two directory levels**: `genre/channel`.
  - Slugify directory names.
- [ ] Update handlers to use this validated path before filesystem or download calls.
- [ ] Remove hard-coded `baseDir` from `downloader` — let it take a full path from caller.
- [ ] Switch downloader to `exec.CommandContext` for cancellation support.
- [ ] Add health endpoints:
  - `/healthz` → always `200 OK`.
  - `/readyz` → checks base dir is writable and `yt-dlp` is present.

## Phase 2 — Middleware & Ergonomics
- [ ] Add request logging middleware (method, path, status, duration).
- [ ] Add request timeout middleware (configurable per-route).
- [ ] Improve error logging (capture `yt-dlp` stderr in downloader).
- [ ] Ensure `.gitignore` excludes build outputs (`bin/`, `data/`) and temp files.
- [ ] Verify Makefile covers: `run`, `test`, `docker-run`, `clean`.

## Phase 3 — Testability
- [ ] Keep `Downloader` interface in `api` (consumer).
- [ ] Implement concrete `YouTube` struct in `downloader` that satisfies the interface.
- [ ] Inject the concrete downloader into server from `main`.
- [ ] Add unit tests for:
  - `GetDirectoriesHandler` (with temp dir).
  - `DownloadHandler` (using mock `Downloader`).

## Phase 4 — Transition to HTMX / SSR
- [ ] Reorganize `internal/api` → `internal/web` for clarity.
- [ ] Add template rendering layer:
  - Layout templates.
  - Page templates.
  - Partial templates for HTMX.
- [ ] Add `/static/` for CSS/JS assets.
- [ ] Keep `/hx/*` routes for HTMX fragment updates (list refresh, create dir, start download).
- [ ] Convert JSON endpoints to HTML/partial responses where appropriate.
- [ ] Serve both pages and partials from the same router (same origin — no CORS needed).

