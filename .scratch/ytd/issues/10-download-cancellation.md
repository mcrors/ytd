Status: done
Phase: 2 — Download Flow

Cancel an in-progress download by cancelling the context passed to the yt-dlp subprocess. After the process exits, delete any `.part` files left in the target directory. Mark the download as `cancelled` in SQLite.

## Acceptance criteria

- [x] Each queued job holds a `context.CancelFunc`; stored so the handler can invoke it by download ID
- [x] `DELETE /downloads/{id}/cancel` HTMX route calls the cancel func
- [x] After subprocess exits on cancellation, glob `*.part` files in target dir and delete them
- [x] Download status set to `cancelled` in SQLite
- [x] Cancelling a completed or already-cancelled download is a no-op (not an error)
- [x] Unit tests for .part file cleanup logic
