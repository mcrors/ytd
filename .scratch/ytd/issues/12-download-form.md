Status: ready-for-agent
Phase: 3 — UI

The main download submission form. URL input, format picker (3 presets), optional rename field, and folder browser trigger. Submits via HTMX POST.

## Acceptance criteria

- [ ] URL input field (required)
- [ ] Format picker: dropdown or radio with 3 options — Best quality (default), 1080p, Audio only (MP3)
- [ ] Optional rename field (leave blank to use YouTube title)
- [ ] Folder browser button opens the folder picker (see ticket 13)
- [ ] Selected folder shown as text below the button
- [ ] Form submits via `hx-post="/downloads"`, swaps response into the queue section
- [ ] Validation: empty URL shows inline error without page reload
- [ ] Mobile-friendly layout
