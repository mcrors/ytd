Status: ready-for-agent
Phase: 4 — Subscriptions

UI for adding, viewing, pausing, resuming, and deleting subscriptions. Main subscriptions page shows status at a glance — no error noise front and centre.

## Acceptance criteria

- [ ] `GET /subscriptions` — full page with list of subscriptions
- [ ] Each row: channel URL, target folder, status badge (active/paused), last checked, last downloaded
- [ ] Add subscription form: channel URL + folder browser (reuses ticket 13 fragment) + submit button
- [ ] Pause/resume toggle per subscription (HTMX POST, updates status in place)
- [ ] Delete button per subscription (with confirmation) — removes subscription and its poll_failures
- [ ] No poll failure details on this page (link to poll history view instead)
- [ ] Empty state when no subscriptions
