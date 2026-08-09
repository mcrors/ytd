# Subscriptions only download videos published after the subscription was created

When a Subscription is created we record the timestamp and pass `--dateafter <created_at>` to yt-dlp on every poll. This means historical videos are never automatically downloaded.

We considered downloading the full back-catalogue on first poll, but rejected it: subscribing to an active channel could silently queue hundreds of videos. The primary use case is "get new content automatically" — historical content is a deliberate one-off Download. A future opt-in backfill flag is noted in the PRD but out of scope for v1.

yt-dlp's `--download-archive` flag is also used to prevent re-downloading within the post-subscription window (e.g. if the app crashes mid-poll and re-runs the same date window).
