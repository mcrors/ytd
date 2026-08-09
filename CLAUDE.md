## Working philosophy

This is a personal portfolio project built AI-first. The goal is for Rory to:

- Fully understand the design and code — not just ship it
- Run high-level and low-level design decisions through AI, since each has blind spots the other doesn't
- Have AI write most of the implementation, but stay involved in patterns, approach, and architecture
- End up with code and docs he can explain to someone a year from now
- Showcase the ability to build systems *with* AI — not vibe-code. Future employers should see a disciplined process: grilled design docs, ADRs, TDD, code review.

**Implications for agents**: never just implement — explain the approach and the trade-offs first. When there's a design decision, surface it. When a pattern is chosen, say why. Write code that Rory can read and own.

---

## Agent skills

### Issue tracker

Issues live as local markdown files under `.scratch/<feature>/`. See `docs/agents/issue-tracker.md`.

### Triage labels

Default canonical labels (needs-triage, needs-info, ready-for-agent, ready-for-human, wontfix). See `docs/agents/triage-labels.md`.

### Domain docs

Single-context layout — one `CONTEXT.md` at repo root plus `docs/adr/`. See `docs/agents/domain.md`.
