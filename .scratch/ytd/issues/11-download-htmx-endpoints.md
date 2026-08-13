Status: done
Phase: 2 — Download Flow

HTMX routes for the full download lifecycle: submit, poll for status, and cancel. All return HTML fragments for HTMX to swap in.

## Routes

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/downloads` | Submit a download; returns a queue-row fragment |
| `GET` | `/downloads/{id}/status` | Returns updated progress fragment (polled every 2s) |
| `POST` | `/downloads/{id}/cancel` | Cancels the download; returns updated row fragment |
| `GET` | `/downloads/history` | Returns the full history fragment |

## Acceptance criteria

- [x] All routes registered in `web.RegisterRoutes()`
- [x] Submit handler: validates URL + format, calls `GetTitle`, inserts row, enqueues job, returns fragment
- [x] Status handler: reads download row from SQLite, renders progress fragment
- [x] Cancel handler: triggers cancellation, returns updated row fragment
- [x] History handler: returns all non-active rows ordered by `created_at` desc
- [x] Integration tests for each route (in-memory SQLite, mock downloader)
