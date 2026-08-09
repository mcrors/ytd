Status: resolved
Phase: 1 — Foundation

SQLite database connection, migration system, and initial schema. Pure Go (`modernc.org/sqlite`, no CGO). Migrations embedded in binary via `//go:embed`.

## Acceptance criteria

- [x] `internal/db/db.go` — `Open()` with WAL mode, foreign keys, busy timeout pragmas
- [x] `internal/db/migrate.go` — `Migrate()` applies numbered SQL files, idempotent
- [x] `internal/db/migrations/001_create_downloads.sql` — downloads table with status, format, error_message columns
- [x] Wired into `main.go` at startup
- [x] Unit tests
