Status: ready-for-agent
Phase: 5 — Library Browser

Read-only file browser for the MEDIA_DIR. Left panel: navigable folder tree. Right panel: file list for the selected folder. No move, rename, or delete.

## Acceptance criteria

- [ ] `GET /library` — full page with two-panel layout
- [ ] `GET /fragments/library/tree` — folder tree fragment (reuses filesystem traversal from ticket 13)
- [ ] `GET /fragments/library/files?path=<rel>` — file list fragment for a selected folder
- [ ] File list shows: filename, size (human-readable), date added
- [ ] Clicking a folder in the tree loads its file list via HTMX (no page reload)
- [ ] No action buttons (read-only)
- [ ] Empty folder state shown
