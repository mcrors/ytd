package downloader_test

import (
	"testing"

	"github.com/mcrors/ytd/internal/downloader"
)

func TestParseProgress(t *testing.T) {
	tests := []struct {
		line    string
		wantPct int
		wantOk  bool
	}{
		// Standard yt-dlp progress lines
		{"[download]  45.3% of 123.45MiB at 2.34MiB/s ETA 00:45", 45, true},
		{"[download]   0.0% of 500.00MiB at 1.00MiB/s ETA 08:20", 0, true},
		{"[download] 100% of 123.45MiB at 5.00MiB/s ETA 00:00", 100, true},
		{"[download]   9.8% of   1.23GiB at  3.45MiB/s ETA 05:12", 9, true},
		// Non-progress lines
		{"[download] Destination: /tmp/video.mp4", 0, false},
		{"[youtube] Extracting URL: https://example.com", 0, false},
		{"[info] Writing video description", 0, false},
		{"", 0, false},
		{"45.3% of 123MiB", 0, false}, // missing [download] prefix
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			pct, ok := downloader.ParseProgress(tt.line)
			if ok != tt.wantOk {
				t.Fatalf("ParseProgress(%q) ok = %v, want %v", tt.line, ok, tt.wantOk)
			}
			if ok && pct != tt.wantPct {
				t.Errorf("ParseProgress(%q) pct = %d, want %d", tt.line, pct, tt.wantPct)
			}
		})
	}
}
