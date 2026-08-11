package download

import "fmt"

type Format string

const (
	FormatBest  Format = "best"
	Format1080p Format = "1080p"
	FormatAudio Format = "audio"
)

// FormatArgs maps a Format preset to the yt-dlp flags it requires.
func FormatArgs(f Format) ([]string, error) {
	switch f {
	case FormatBest:
		return []string{"-f", "bestvideo+bestaudio/best"}, nil
	case Format1080p:
		return []string{"-f", "bestvideo[height<=1080]+bestaudio/best[height<=1080]"}, nil
	case FormatAudio:
		return []string{"-f", "bestaudio/best", "--extract-audio", "--audio-format", "mp3"}, nil
	default:
		return nil, fmt.Errorf("unknown format %q: must be one of best, 1080p, audio", f)
	}
}
