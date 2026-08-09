Status: ready-for-agent
Phase: 3 — UI

DaisyUI toast notifications for download completion, failure, and cancellation. Triggered via HTMX response headers (`HX-Trigger`) — no polling needed.

## Acceptance criteria

- [ ] Toast container div in `layout.html`, positioned fixed bottom-right
- [ ] Server sets `HX-Trigger: {"showToast": {"message": "...", "type": "success|error"}}` on relevant responses
- [ ] HTMX event listener renders a DaisyUI alert into the toast container
- [ ] Toast auto-dismisses after 4s
- [ ] Success toast on download completion
- [ ] Error toast on download failure (with short error summary)
- [ ] Info toast on cancellation
