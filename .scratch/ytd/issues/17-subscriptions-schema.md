Status: ready-for-agent
Phase: 4 — Subscriptions

SQLite migration for the subscriptions table and poll_failures table.

## Schema

```sql
CREATE TABLE subscriptions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    channel_url TEXT NOT NULL,
    target_dir TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',   -- active | paused
    archive_path TEXT NOT NULL DEFAULT '',   -- path to yt-dlp --download-archive file
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_checked_at DATETIME,
    last_downloaded_at DATETIME
);

CREATE TABLE poll_failures (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    subscription_id INTEGER NOT NULL REFERENCES subscriptions(id),
    error_message TEXT NOT NULL,
    failed_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## Acceptance criteria

- [ ] `internal/db/migrations/002_create_subscriptions.sql` created with above schema
- [ ] `status` constrained to `active` | `paused`
- [ ] Migration runs cleanly on top of 001; idempotent on re-run
- [ ] `archive_path` defaults to empty string — set to a file path alongside the target_dir when first poll runs
