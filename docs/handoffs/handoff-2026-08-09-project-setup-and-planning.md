# Handoff — 2026-08-09 — Project setup and planning

## What happened this session

This was a design and planning session, not an implementation session. No production code was written or changed.

1. **Read the repo** — existing code (Go backend, yt-dlp subprocess wrapper, SQLite, gorilla/mux) and all docs were reviewed. Phase 1.1 (config) and 1.2 (SQLite) were already done.

2. **Set up agent skills** — `CLAUDE.md`, `docs/agents/issue-tracker.md`, `docs/agents/domain.md`, `docs/agents/triage-labels.md` created. Issue tracker is local markdown (`.scratch/`).

3. **Working philosophy captured** — `~/.claude/CLAUDE.md` (global, loads in every session) and the `## Working philosophy` section of `CLAUDE.md` (project-local) both describe how Rory wants to work: AI writes the code, but Rory stays in the loop on design decisions and must be able to explain everything a year from now. Portfolio-grade process matters.

4. **Grilled the PRD** — 11 design decisions resolved. See `CONTEXT.md`, `docs/adr/0001-htmx-only-no-json-api.md`, `docs/adr/0002-subscription-date-cutoff.md`, and the updated `docs/PRD.md` for the outcomes. Key decisions:
   - `internal/api` package removed entirely — HTMX-only transport
   - In-memory download queue (not persisted across restarts)
   - 2s HTMX polling for progress (not SSE)
   - 3 curated format presets: Best, 1080p, Audio only
   - Title pre-fetched with `--get-title` before queuing
   - Cancellation cleans up `.part` files
   - Subscriptions: `--dateafter <created_at>` + `--download-archive` for dedup
   - Single global `POLL_INTERVAL` (not per-subscription)
   - Poll failures logged separately, not front-and-centre
   - Library browser is read-only in v1

5. **Created 23 tickets** — all in `.scratch/ytd/issues/`. Tickets 01 and 02 are already marked DONE. The project map is at `.scratch/ytd/map.md`.

## Current state

- All planning artefacts written, nothing to implement yet.
- **Next ticket: `03-remove-api-package.md`** — the first implementation task.
- Tickets are numbered sequentially; rename to `NN-slug-DONE.md` when complete.

## Key file locations

| Artefact | Path |
|----------|------|
| Working philosophy | `~/.claude/CLAUDE.md` (global) and `CLAUDE.md` (project) |
| Domain glossary | `CONTEXT.md` |
| ADRs | `docs/adr/` |
| Updated PRD | `docs/PRD.md` |
| Project map | `.scratch/ytd/map.md` |
| All tickets | `.scratch/ytd/issues/` |
| Agent skills config | `docs/agents/` |

## What the next agent should do

Start ticket `03-remove-api-package.md`. Read the ticket file for acceptance criteria. Before writing any code, explain the approach and the trade-offs — per the working philosophy in `CLAUDE.md`.

Tickets 03, 04, and 05 are the remaining Phase 1 (Foundation) tasks and can largely be worked in order. They are prerequisites for all of Phase 2.

## Suggested skills

- `/tdd` — Rory wants to understand the code; test-first keeps the design honest and gives him executable specs to read
- `/code-review` — run after each ticket is implemented
- `/domain-modeling` — if new domain terms emerge during implementation, update `CONTEXT.md` inline
