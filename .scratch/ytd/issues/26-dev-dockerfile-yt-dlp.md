Status: ready-for-agent
Phase: 4 — Hardening

Install yt-dlp in `Dockerfile.dev` so the dev environment can actually run downloads end-to-end without yt-dlp missing from PATH.

## Background

`Dockerfile.dev` is based on `debian:bookworm-slim` and only installs `curl`. yt-dlp is not present, so all downloads fail immediately with "yt-dlp not found in PATH". The app handles this gracefully but the dev environment is not useful for manual testing of the download flow.

## Acceptance criteria

- [ ] `Dockerfile.dev` installs yt-dlp (via pip3 or the official GitHub release binary)
- [ ] `make dev-up` produces a container where a real YouTube download completes successfully
