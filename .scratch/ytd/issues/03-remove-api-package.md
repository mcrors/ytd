Status: done
Phase: 1 — Foundation

Remove the `internal/api` package entirely and replace with `internal/web` using stdlib `http.ServeMux` (Go 1.22+). Remove `gorilla/mux` from `go.mod`. This is the foundation for HTMX-only transport — see ADR-0001.

## Acceptance criteria

- [x] `internal/api/` deleted
- [x] `gorilla/mux` removed from `go.mod` and `go.sum`
- [x] `internal/web/` created with `server.go`, `handlers.go`, `respond.go`
- [x] `web.RegisterRoutes(mux *http.ServeMux, ...)` pattern — registers routes on a shared mux
- [x] `main.go` updated: single `http.ServeMux`, web routes mounted
- [x] `GET /healthz` and `GET /readyz` still work (move them into web package)
- [x] App boots and serves requests
