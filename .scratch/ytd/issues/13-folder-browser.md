Status: ready-for-agent
Phase: 3 — UI

Navigable tree of the MEDIA_DIR filesystem. Used to select a target folder for downloads and subscriptions. Also supports creating new folders at arbitrary depth.

## Acceptance criteria

- [ ] `GET /fragments/folders` returns an HTMX fragment with the folder tree rooted at MEDIA_DIR
- [ ] Clicking a folder selects it and highlights it; selection is posted back to the parent form
- [ ] `POST /fragments/folders` creates a new folder (uses `SafeJoin`); re-renders the tree
- [ ] Folders read from filesystem on demand — not cached
- [ ] Traversal attempts (e.g. `../../etc`) return 400
- [ ] New folder name input with a create button, inline in the tree
- [ ] Works inside a DaisyUI modal triggered from the download form
