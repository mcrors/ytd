Status: ready-for-agent
Phase: 5 — Library Browser

File metadata sourcing for the library browser. SQLite is the primary source for files downloaded via ytd; filesystem `stat` is the fallback for files that exist on disk but have no database record.

## Metadata sources

| Field | Source (SQLite) | Fallback (filesystem) |
|-------|----------------|----------------------|
| Title | `downloads.title` | Filename without extension |
| Size | `stat` (always) | `stat` |
| Date added | `downloads.created_at` | `stat.ModTime()` |

## Acceptance criteria

- [ ] `internal/library/metadata.go` — `ResolveMetadata(path string, db *sql.DB) FileMetadata`
- [ ] Looks up `downloads` row by matching `filename` column; falls back to stat if not found
- [ ] Returns consistent `FileMetadata` struct regardless of source
- [ ] Unit tests for both paths (SQLite hit and miss)
