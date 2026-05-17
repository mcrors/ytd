# YT-DLP Homelab Downloader — Product Requirements Document

## Overview

A self-hosted web application for downloading YouTube videos to an NFS share. Built for personal homelab use, prioritising simplicity and maintainability over scale. The app provides a clean UI for one-off downloads, folder organisation, and channel subscriptions.

---

## Tech Stack

| Layer | Choice | Notes |
|---|---|---|
| Backend | Go | Idiomatic, minimal abstraction |
| Frontend | HTMX + DaisyUI (CDN) | No build step, no Node.js |
| Downloader | yt-dlp | Called as a subprocess |
| Database | SQLite | App state only |
| Container | Docker | Single image deployment |
| Orchestration | k3s | Kubernetes manifests provided separately |

---

## Configuration

Provided via environment variables or a YAML config file mounted into the container.

| Key | Description | Default |
|---|---|---|
| `NFS_MOUNT_PATH` | Absolute path to the NFS share mount | `/mnt/media` |
| `POLL_INTERVAL` | How often to check subscribed channels | `1h` |
| `MAX_CONCURRENT_DOWNLOADS` | Max parallel yt-dlp processes | `2` |
| `DB_PATH` | Path to SQLite database file | `/data/ytdlp.db` |
| `PORT` | HTTP server port | `8080` |

---

## Features

### 1. One-Off Download

- User pastes a YouTube URL into the UI
- User selects a target folder from the folder browser
- Download is queued and executed via yt-dlp
- User can select download format (video quality or audio-only MP3)
- User can change the name of the downloaded video
- Progress is shown live in the UI via HTMX polling

### 2. Folder Management

- User can browse the folder structure of the NFS share from the UI
- User can create new folders and nested subfolders (e.g. `history/mary-beard`)
- Folder tree is read directly from the filesystem on demand — not cached
- No deletion of folders from the UI (to avoid accidents)

### 3. Download Queue & History

- Active downloads displayed with live progress bar and status badge
- Completed and failed downloads retained in history view
- Each history entry shows: title, target folder, format, timestamp, status
- Failed downloads show error reason
- User can cancel an in-progress download

### 4. Channel Subscriptions

- User can subscribe to a YouTube channel by providing its URL
- Each subscription has a configured target folder and check interval
- Background worker polls subscribed channels at the configured interval
- New videos are automatically queued for download
- User can view, pause, and delete subscriptions
- Last checked and last downloaded timestamps shown per subscription

### 5. Library Browser

- Displays files already present on the NFS share
- Navigable folder tree on the left, file list on the right
- File metadata shown: title, size, date added
- Metadata sourced from SQLite where available, filesystem otherwise
- No playback — browse and organise only

---

## UI

- DaisyUI components throughout — cards, badges, progress bars, toasts, modals
- HTMX for all dynamic interactions — no page reloads
- Go `html/template` for server-side rendering
- Responsive layout suitable for desktop, tablet, and mobile
- Mobile-first approach — usable from a phone browser for quick one-off downloads
- Toast notifications for job completion, errors, and subscription events

---

## Out of Scope (v1)

- User authentication
- Multi-user support
- Video playback in the browser
- Non-YouTube sources

---

## Future Considerations

- Webhook or notification (e.g. Gotify, ntfy) on subscription download
- Metadata embedding (thumbnail, description saved alongside video file)
- Scheduled downloads (download at a specific time)
- Import/export of subscription list
- Mobile-optimised native-style layout
