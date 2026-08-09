# HTMX-only transport — no JSON API

This is a personal homelab tool consumed exclusively from a browser. We considered maintaining a full JSON API (`/api/*`) alongside HTMX routes as a portfolio signal, but rejected it: the portfolio goal is demonstrating disciplined process, not API design. Full parity would double the surface area and tests for every feature, with no current consumer beyond the browser. The existing `/api/*` package is removed entirely; all routes serve HTMX fragments. A JSON API can be added later if a real consumer appears.
