Status: done
Phase: 2 — Download Flow

Three curated download format presets. The user picks a label; the service layer maps it to the yt-dlp format string. The raw format string is never exposed in the UI or API.

## Presets

| Label | yt-dlp format string |
|-------|----------------------|
| `best` | `bestvideo+bestaudio/best` |
| `1080p` | `bestvideo[height<=1080]+bestaudio/best[height<=1080]` |
| `audio` | `bestaudio/best` with `--extract-audio --audio-format mp3` |

## Acceptance criteria

- [x] `Format` type defined (e.g. string enum: `best`, `1080p`, `audio`)
- [x] `FormatArgs(f Format) []string` maps preset to yt-dlp flags
- [x] `Downloader.Download()` accepts a `Format` parameter and passes the correct flags
- [x] Invalid format value rejected with a clear error
- [x] Unit tests for each preset mapping
