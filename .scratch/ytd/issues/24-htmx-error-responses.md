Status: ready-for-agent
Phase: 3 — UI

All handler errors currently respond with JSON via `respondError`. In an HTMX context this is wrong — HTMX will swap raw JSON into the target element. Errors need to be surfaced as HTML, either via a toast trigger or an inline error fragment.

## Background

`respondError` in `internal/web/respond.go` writes `application/json`. This was left over from the original JSON API design. Now that all routes return HTML fragments, error responses should follow suit.

## Options to consider

- Return an error fragment that HTMX can swap in (shows inline error where the action was)
- Use `HX-Trigger: showToast` header (ties into ticket 16) and return 4xx/5xx with no body
- Hybrid: inline fragment for submit errors, toast for unexpected server errors

## Acceptance criteria

- [ ] `respondError` replaced or supplemented with an HTMX-compatible error path
- [ ] Validation errors (400) surfaced visibly to the user in the UI
- [ ] Server errors (500) surfaced as a toast or fallback message
- [ ] No raw JSON reaches the DOM in normal HTMX flows
