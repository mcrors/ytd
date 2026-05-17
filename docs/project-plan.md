# YT-DLP Homelab Downloader — Project Plan

## Current State

- **Go backend**: HTTP server with gorilla/mux, JSON API
- **yt-dlp wrapper**: `Download` and `GetChannel` via subprocess, testable interfaces
- **Endpoints**: `POST /api/download`, `GET /api/directories`, `POST /api/directory`, health/readiness probes
- **Folder logic**: two-level enforced (`genre/channel`), slugified — to be replaced with flexible user-defined paths
- **Frontend**: none
- **Database**: none
- **Queue**: none — downloads are synchronous in the request handler
- **Config**: `YTD_BASE_DIR` env var only
- **Docker**: Dockerfile and docker-compose.yaml exist but not validated against current code

---

## Phase 1: Foundation

**Goal**: infrastructure layer that everything else builds on.

- SQLite setup — schema, migrations, connection management
- Config system — all env vars loaded and validated at startup (`MEDIA_DIR`, `DB_PATH`, `PORT`, `MAX_CONCURRENT_DOWNLOADS`, `POLL_INTERVAL`)
- Base HTML layout template — DaisyUI + HTMX via CDN, Go `html/template`
- Static file serving for local assets
- Dual-transport routing — existing JSON API under `/api/*`, new HTMX fragment routes alongside
- Path handling refactor — remove `normalizeTwoLevel`, replace with flexible paths + traversal protection

**Ends when**: app boots, serves a base HTML page, connects to SQLite, reads config from env, JSON API still works.

---

## Phase 2: Download Flow

**Goal**: core download feature end-to-end.

- Download queue — channel-based worker pool, bounded by `MAX_CONCURRENT_DOWNLOADS`
- yt-dlp progress parsing from stdout (percentage, speed, ETA)
- Format selection — video quality presets + audio-only MP3
- Download cancellation via context
- Flexible folder targeting — user picks full path under `MEDIA_DIR`
- Download history persisted to SQLite (title, folder, format, timestamp, status, error)
- Both JSON and HTMX endpoints for submitting, querying status, cancelling

**Ends when**: you can submit a URL, pick a folder, pick a format, watch progress, cancel, and see history — via both API and HTML.

---

## Phase 3: UI

**Goal**: fully usable browser interface for one-off downloads.

- Download form with URL input, format selector, folder browser
- Folder browser — tree navigation of `MEDIA_DIR`, read from filesystem on demand
- Folder creation from the UI (arbitrary depth)
- Active download queue — live progress bars via HTMX polling (`hx-trigger="every 2s"`)
- Download history view with status badges
- Toast notifications for completion, failure, errors
- Mobile-first responsive layout

**Ends when**: the app is fully usable from a phone browser for one-off downloads.

---

## Phase 4: Subscriptions

**Goal**: subscribe to channels, auto-download new videos.

- Subscription model in SQLite (channel URL, target folder, check interval, status, timestamps)
- Background poller goroutine — ticks on `POLL_INTERVAL`, checks each active subscription
- New videos auto-queued through the existing download queue
- Subscription management UI — add, pause, resume, delete
- Per-subscription status display — last checked, last downloaded, error state
- Both JSON and HTMX endpoints

**Ends when**: you can subscribe to a channel and new videos download automatically to the configured folder.

---

## Phase 5: Library Browser

**Goal**: browse what's already on disk.

- Folder tree navigation (left panel)
- File list for selected folder (right panel)
- File metadata — title, size, date added
- Metadata from SQLite where available, filesystem `stat` otherwise
- No playback — browse and organise only

**Ends when**: you can navigate the full `MEDIA_DIR` tree and see what's there.

---

## Phase 6: Deployment

**Goal**: production-ready container and k3s deployment.

- Dockerfile — multi-stage build, yt-dlp baked in
- docker-compose.yaml for local development
- k3s manifests — Deployment, Service, PVC for `MEDIA_DIR`, PVC for SQLite, ConfigMap for env vars
- Health/readiness probes wired to existing `/healthz` and `/readyz`

**Ends when**: `docker compose up` works locally, k3s manifests deploy to the cluster.

---

## Dependency Graph

```
Phase 1 (foundation) ──→ Phase 2 (download flow) ──→ Phase 3 (UI)
                                                  ──→ Phase 4 (subscriptions)
                                                  ──→ Phase 5 (library browser)

Phase 6 (deployment) can proceed in parallel from Phase 2 onward.
Phase 4 and Phase 5 are independent of each other.
```

---

## Open Decisions

| Item | Options | Notes |
|------|---------|-------|
| JSON API scope | Full parity with HTMX / download-only / drop later | Currently leaning: full parity as portfolio signal |
| Folder depth limit | Unlimited / cap at N levels | Leaning unlimited with traversal protection |
| Config file support | Env vars only / env + YAML | PRD says both; env-only is simpler to start |
| Auto channel subfolder | User picks full path / option to auto-append channel name | Dropped from service layer; revisit if wanted |
| Download format presets | yt-dlp format strings exposed raw / curated preset list | Curated is friendlier; raw is more flexible |
| SQLite migrations | Embedded in Go binary / separate migration files | Embedded (via `embed` package) is simpler for single-binary deployment |
| Progress polling interval | Fixed 2s / configurable / SSE instead of polling | 2s polling is simplest; SSE is smoother but more complex |
