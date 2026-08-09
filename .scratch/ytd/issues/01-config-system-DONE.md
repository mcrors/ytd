Status: resolved
Phase: 1 — Foundation

Centralised config loaded from YAML file + env vars with sensible defaults. Priority chain: env vars > config file > defaults. All fields prefixed `YTD_`.

## Acceptance criteria

- [x] `internal/config/config.go` with `Config` struct and `Load()` function
- [x] YAML config file support (`config.yaml`)
- [x] Env var overrides: `YTD_MEDIA_DIR`, `YTD_DB_PATH`, `YTD_PORT`, `YTD_MAX_CONCURRENT_DOWNLOADS`, `YTD_POLL_INTERVAL`
- [x] App starts with no config file and no env vars (all defaults)
- [x] Unit tests
