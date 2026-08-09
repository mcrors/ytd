Status: ready-for-agent
Phase: 1 — Foundation

Replace the rigid `normalizeTwoLevel`/`slugify` path logic with a flexible `SafeJoin` that allows arbitrary depth paths while blocking traversal attacks.

## Acceptance criteria

- [ ] `normalizeTwoLevel()`, `slugify()`, and `twoSegRe` removed from `internal/api/filesystem.go` (or equivalent location after ticket 03)
- [ ] `SafeJoin(baseDir, userPath string) (string, error)` implemented: joins, cleans, verifies result is under baseDir
- [ ] `download/service.go` uses `SafeJoin`; `GetChannel` call and auto-channel-name injection removed from download flow
- [ ] Directory creation handler uses `SafeJoin`
- [ ] Unit tests: normal path, traversal blocked, absolute path blocked, empty path returns base
- [ ] Traversal attempts return HTTP 400
