Status: ready-for-agent
Phase: 4 — Subscriptions

Separate page (or expandable section per subscription) showing the poll failure log. Not surfaced on the main subscriptions page — opt-in for when you want to investigate.

## Acceptance criteria

- [ ] `GET /subscriptions/{id}/poll-history` returns a page or fragment showing `poll_failures` for that subscription
- [ ] Each row: timestamp, error message
- [ ] Ordered by `failed_at` desc; cap at 100 rows
- [ ] Linked from the subscription row (e.g. "N failures" badge/link)
- [ ] Empty state when no failures
