Status: ready-for-agent
Phase: 3 — UI

Live view of in-progress downloads. Each row shows title, target folder, format, progress bar, and a cancel button. Refreshes every 2s via HTMX polling.

## Acceptance criteria

- [ ] `GET /fragments/queue` returns the active queue fragment (status = `queued` or `downloading`)
- [ ] Each row: title, folder, format badge, DaisyUI progress bar (driven by `downloads.progress`), cancel button
- [ ] Fragment polls itself every 2s via `hx-trigger="every 2s"` while any active downloads exist
- [ ] Polling stops automatically when queue is empty (swap out the trigger)
- [ ] Cancel button posts to `/downloads/{id}/cancel`; row updates in place
- [ ] Empty state shown when queue is clear
