Status: ready-for-agent
Phase: 4 — Hardening

Replace `log.Printf` calls throughout the codebase with `log/slog` (Go standard structured logger). Add log levels so debug noise can be suppressed in production.

## Background

Currently all logging uses `log.Printf` with no level distinction. In production, queue progress updates and routine request logs shouldn't appear at the same level as errors. `log/slog` is the idiomatic Go 1.21+ choice — no external dependency, structured output, level filtering.

## Acceptance criteria

- [ ] `log/slog` configured at startup in `main.go` with level from config/env (`YTD_LOG_LEVEL`, default `info`)
- [ ] All `log.Printf` calls replaced with appropriate `slog` level (`slog.Info`, `slog.Warn`, `slog.Error`, `slog.Debug`)
- [ ] Queue progress updates (`download %d completed/failed/cancelled`) at `Info`
- [ ] DB exec errors in queue worker at `Warn` (non-fatal)
- [ ] Request middleware logging at `Info`
- [ ] Structured fields where useful (e.g. `slog.Int64("id", job.ID)`)
