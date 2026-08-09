# ytd

A self-hosted homelab app for downloading YouTube videos to an NFS share. Single user, browser-only interface, deployed on k3s.

## Language

**Download**:
A request to fetch a single YouTube video via yt-dlp and save it to a target directory under the media root.
_Avoid_: job, task, transfer

**Queue**:
The in-memory worker pool that executes Downloads concurrently, bounded by `MAX_CONCURRENT_DOWNLOADS`. Not persisted — cleared on restart.
_Avoid_: job queue, task queue

**Format**:
A curated preset for the yt-dlp output: `best` (best available quality), `1080p`, or `audio` (MP3). Maps to a yt-dlp format string in the service layer. The user picks a Format; the raw yt-dlp string is never exposed.
_Avoid_: quality, resolution, yt-dlp format string

**Subscription**:
A configured channel URL + target directory that the background poller checks on a schedule, automatically queuing new videos for Download. Only videos published after the Subscription was created are ever queued — historical content is excluded.
_Avoid_: feed, channel watch, auto-download

**Media root**:
The base directory on the NFS share (`MEDIA_DIR`) that all Downloads and Subscriptions write into. User-facing paths are always relative to this root.
_Avoid_: base dir, NFS mount, download directory
