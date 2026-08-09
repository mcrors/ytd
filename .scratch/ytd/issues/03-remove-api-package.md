Status: ready-for-agent
Phase: 1 — Foundation

Remove the `internal/api` package entirely and replace with `internal/web` using stdlib `http.ServeMux` (Go 1.22+). Remove `gorilla/mux` from `go.mod`. This is the foundation for HTMX-only transport — see ADR-0001.

## Acceptance criteria

- [ ] `internal/api/` deleted
- [ ] `gorilla/mux` removed from `go.mod` and `go.sum`
- [ ] `internal/web/` created with `server.go`, `handlers.go`, `respond.go`
- [ ] `web.RegisterRoutes(mux *http.ServeMux, ...)` pattern — registers routes on a shared mux
- [ ] `main.go` updated: single `http.ServeMux`, web routes mounted
- [ ] `GET /healthz` and `GET /readyz` still work (move them into web package)
- [ ] App boots and serves requests
