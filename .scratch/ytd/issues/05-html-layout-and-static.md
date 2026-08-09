Status: ready-for-agent
Phase: 1 — Foundation

Embedded Go templates with DaisyUI + HTMX (both via CDN). Base layout, placeholder pages, static file serving. Dev mode reloads templates from disk without recompile.

## Acceptance criteria

- [ ] `templates/layout.html` — base layout with DaisyUI CDN, HTMX CDN, nav bar, `{{ block "content" }}` slot, toast container, responsive viewport meta
- [ ] `templates/pages/index.html` — placeholder landing page
- [ ] Templates embedded via `//go:embed`; loaded with `template.ParseFS()` at startup
- [ ] `static/` directory embedded and served at `GET /static/*`
- [ ] `YTD_DEV=true` reloads templates from disk on each request (no recompile during frontend work)
- [ ] `GET /` returns a styled HTML page with DaisyUI and HTMX loaded
- [ ] `GET /static/style.css` returns a file
