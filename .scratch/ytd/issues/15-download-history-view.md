Status: ready-for-agent
Phase: 3 — UI

History of completed, failed, and cancelled downloads. Each entry shows title, folder, format, timestamp, and status badge. Failed entries show the error reason.

## Acceptance criteria

- [ ] `GET /downloads/history` returns history fragment ordered by `created_at` desc
- [ ] Each row: title, target folder, format, timestamp, DaisyUI status badge (completed/failed/cancelled)
- [ ] Failed rows show `error_message` inline (collapsed by default, expandable)
- [ ] Pagination or a "show more" if history is long (cap initial view at 50 rows)
- [ ] Empty state message when no history
