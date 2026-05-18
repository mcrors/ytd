# Phase 1: Foundation

## Goal

Infrastructure layer that everything else builds on. App boots, serves a base HTML page, connects to SQLite, reads config from env/file, JSON API still works.

---

## 1.1 — Config System - DONE

**Goal**: centralised config loaded from YAML file + env vars with sensible defaults.

Priority chain: **env vars > config file > defaults**.

- [ ] Create `internal/config/config.go`
  - `Config` struct with fields:
    ```go
    type Config struct {
        MediaDir             string        // where downloaded media lives
        DBPath               string        // SQLite database path
        Port                 string        // HTTP server port
        MaxConcurrentDL      int           // download worker pool size
        PollInterval         time.Duration // subscription check frequency
    }
    ```
  - `Load(configPath string) (*Config, error)` — reads YAML if file exists, applies env overrides, fills defaults
- [ ] YAML config file format (`config.yaml`):
  ```yaml
  media_dir: /mnt/media
  db_path: /data/ytdlp.db
  port: "8080"
  max_concurrent_downloads: 2
  poll_interval: 1h
  ```
- [ ] Env var mapping (all prefixed `YTD_`):
  | Env Var | Config Field | Default |
  |---------|-------------|---------|
  | `YTD_MEDIA_DIR` | MediaDir | `./data/media` |
  | `YTD_DB_PATH` | DBPath | `./data/ytdlp.db` |
  | `YTD_PORT` | Port | `8080` |
  | `YTD_MAX_CONCURRENT_DOWNLOADS` | MaxConcurrentDL | `2` |
  | `YTD_POLL_INTERVAL` | PollInterval | `1h` |
  | `YTD_CONFIG_PATH` | — | `./config.yaml` (overrides default config file location) |
- [ ] Remove inline `os.Getenv("YTD_BASE_DIR")` from `main.go`, replace with `config.Load()`
- [ ] Validate: app starts with no config file and no env vars (all defaults)
- [ ] Validate: app starts with config file only
- [ ] Validate: env var overrides config file value

**No external dependencies.** `os`, `encoding/yaml` (stdlib doesn't have YAML — use `gopkg.in/yaml.v3`, it's the standard choice and a single dependency) for file parsing, `os.Getenv` for env.

---

## 1.2 — SQLite

**Goal**: database connection, migration system, initial schema.

- [ ] Add `modernc.org/sqlite` dependency (pure Go, no CGO)
- [ ] Create `internal/db/db.go`
  - `Open(dbPath string) (*sql.DB, error)` — opens connection, sets pragmas:
    ```go
    db.Exec("PRAGMA journal_mode=WAL")
    db.Exec("PRAGMA foreign_keys=ON")
    db.Exec("PRAGMA busy_timeout=5000")
    ```
  - Returns `*sql.DB` — standard library interface, no wrapper
- [ ] Create `internal/db/migrate.go`
  - `migrations/` directory with numbered SQL files, embedded via `//go:embed`
  - `schema_migrations` table tracks applied migrations:
    ```sql
    CREATE TABLE IF NOT EXISTS schema_migrations (
        version INTEGER PRIMARY KEY,
        applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    ```
  - `Migrate(db *sql.DB) error` — reads embedded files, compares against `schema_migrations`, applies new ones in a transaction
  - Migration files named: `001_create_downloads.sql`, `002_create_subscriptions.sql`, etc.
- [ ] Create `internal/db/migrations/001_create_downloads.sql`:
  ```sql
  CREATE TABLE downloads (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      url TEXT NOT NULL,
      title TEXT NOT NULL DEFAULT '',
      target_dir TEXT NOT NULL,
      filename TEXT NOT NULL DEFAULT '',
      format TEXT NOT NULL DEFAULT 'video',
      status TEXT NOT NULL DEFAULT 'queued',
      error_message TEXT NOT NULL DEFAULT '',
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  - `status` values: `queued`, `downloading`, `completed`, `failed`, `cancelled`
  - `format` values: `video`, `audio`
- [ ] Wire into `main.go`: `db.Open(cfg.DBPath)` → `db.Migrate(db)` at startup
- [ ] Validate: app starts, creates `ytdlp.db`, `schema_migrations` table exists, `downloads` table exists
- [ ] Validate: restart app — migration is idempotent, no error

**Note**: subscription table (`002_create_subscriptions.sql`) will be defined in Phase 4. The migration system supports adding files later without touching existing code.

---

## 1.3 — Path Handling Refactor

**Goal**: replace rigid two-level paths with flexible user-defined paths, protected against traversal.

- [ ] Remove `normalizeTwoLevel()` and `slugify()` from `internal/api/filesystem.go`
- [ ] Remove `twoSegRe` regex
- [ ] Create `SafeJoin(baseDir, userPath string) (string, error)` — probably in a new `internal/pathutil/pathutil.go` or keep in `filesystem.go`:
  ```
  1. filepath.Join(baseDir, userPath)
  2. filepath.Clean the result
  3. Verify result starts with filepath.Clean(baseDir)
  4. Reject if not (traversal attempt)
  ```
- [ ] Update `download/service.go` — remove `normalizeTwoLevel` call, use `SafeJoin` instead
- [ ] Remove auto channel name injection (`ds.downloader.GetChannel` + appending to path) — user controls the full path
- [ ] Update `createDirectoryHandler` to use `SafeJoin` instead of `normalizeTwoLevel`
- [ ] Update `getDirectoriesHandler` to return recursive tree (not just top-level entries)
- [ ] Write tests for `SafeJoin`:
  - Normal path: `SafeJoin("/mnt/media", "History/Mary Beard")` → `/mnt/media/History/Mary Beard`
  - Traversal blocked: `SafeJoin("/mnt/media", "../../etc/passwd")` → error
  - Absolute path blocked: `SafeJoin("/mnt/media", "/etc/passwd")` → error
  - Empty path: `SafeJoin("/mnt/media", "")` → `/mnt/media` (valid, returns base)
- [ ] Validate: existing directory creation still works with new paths
- [ ] Validate: traversal attempts return 400

---

## 1.4 — Base HTML Layout

**Goal**: embedded Go templates with DaisyUI + HTMX, serving a base page.

- [ ] Create `templates/` directory at project root:
  ```
  templates/
  ├── layout.html          # base layout: <html>, <head>, nav, content block, toast container
  ├── pages/
  │   └── index.html       # landing page (placeholder for now)
  └── partials/
      └── toast.html        # toast notification fragment
  ```
- [ ] `layout.html` includes:
  - DaisyUI CSS via CDN (`<link>` tag)
  - HTMX JS via CDN (`<script>` tag)
  - `<nav>` bar with app name + page links (Downloads, Library, Subscriptions — placeholder hrefs)
  - `{{ block "content" . }}{{ end }}` for page body
  - Toast container div (empty, HTMX swaps notifications into it)
  - Responsive meta viewport tag
- [ ] Embed templates via `//go:embed templates/*` in the web package
- [ ] `template.ParseFS()` to load all templates at startup
- [ ] Validate: `GET /` returns a rendered HTML page with DaisyUI styling and HTMX loaded
- [ ] Validate: page is responsive on mobile viewport

---

## 1.5 — Dual-Transport Routing

**Goal**: JSON API and HTML handlers coexist, both calling the same service layer. Single server, single `ServeMux`, two handler packages.

- [ ] Create `internal/web/` package for HTMX/HTML handlers:
  ```
  internal/web/
  ├── server.go           # RegisterRoutes func, template loading
  ├── handlers.go         # HTML fragment handlers
  └── respond.go          # respondHTML helper
  ```
- [ ] `internal/api/` stays as-is — JSON transport under `/api/*`
  - [ ] Refactor `api.NewServer()` to `api.RegisterRoutes(mux, ...)` — registers routes on a shared mux instead of returning its own handler
  - [ ] Remove gorilla/mux dependency, use stdlib `http.ServeMux` (Go 1.22+ method routing)
- [ ] `internal/web/` serves:
  - `GET /` — index page (full page render)
  - `GET /downloads` — downloads page (full page render)
  - `GET /library` — library page (full page, Phase 5)
  - `GET /subscriptions` — subscriptions page (full page, Phase 4)
  - Fragment endpoints as needed (e.g. `GET /fragments/directories` for folder tree partial)
- [ ] `web.RegisterRoutes(mux, ...)` — same pattern as api
- [ ] Update `main.go`:
  ```go
  mux := http.NewServeMux()
  api.RegisterRoutes(mux, downloadService)   // mounts /api/*
  web.RegisterRoutes(mux, downloadService)   // mounts /, /downloads, etc.

  srv := &http.Server{Addr: ":" + cfg.Port, Handler: mux}
  ```
- [ ] Both packages receive the same service interfaces — no duplication of business logic
- [ ] Dev mode template reload: `YTD_DEV=true` env var causes templates to be read from disk on each request instead of from embedded FS. Faster iteration during frontend work, no recompile needed for template changes.
- [ ] Validate: `curl /api/directories` returns JSON
- [ ] Validate: `GET /` in browser returns styled HTML page
- [ ] Validate: gorilla/mux fully removed from `go.mod`

---

## 1.6 — Static File Serving

**Goal**: serve local static assets (CSS overrides, favicon, any JS if needed).

- [ ] Create `static/` directory at project root
- [ ] Embed via `//go:embed static/*`
- [ ] Mount at `GET /static/*` using `http.FileServer`
- [ ] Minimal contents for now — a `favicon.ico` and a `style.css` for any DaisyUI overrides
- [ ] Validate: `GET /static/style.css` returns the file

---

## Dependency Graph

```
1.1 (config) ──→ 1.2 (SQLite) ──→ main.go wiring
                                      ↑
1.3 (path refactor) ─────────────────┘
                                      ↑
1.4 (templates) ──→ 1.5 (routing) ───┘
                                      ↑
1.6 (static files) ──────────────────┘
```

1.1 and 1.3 can proceed in parallel.
1.4 and 1.6 can proceed in parallel.
1.2 depends on 1.1 (needs `DBPath` from config).
1.5 depends on 1.4 (needs templates) and integrates everything.

---

## Decisions Made

| Item | Decision | Notes |
|------|----------|-------|
| Router | stdlib `http.ServeMux` (Go 1.22+) | Drop gorilla/mux. Dependency-free, method routing built in. |
| Config file path | `./config.yaml` default, override via `YTD_CONFIG_PATH` | Zero config for local dev, mounted YAML for k3s. |
| Env var prefix | `YTD_` | All env vars prefixed to avoid collisions. |
| GetChannel | Keep in downloader, remove from download flow | Available for future use (subscriptions, auto-naming). Not called during one-off downloads. |
| Dev mode | `YTD_DEV=true` reloads templates from disk | No recompile needed for template changes during frontend work. |

---

## Open Decisions

| Item | Options | Notes |
|------|---------|-------|
| YAML library | `gopkg.in/yaml.v3` | Only external dep for config. Standard choice. |
| Download format column | `video` / `audio` enum / free string | Enum is simpler now. Free string allows future presets (720p, 1080p, 4k). Decide in Phase 2. |
