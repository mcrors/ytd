package downloader

import (
	"regexp"
	"strconv"
	"strings"
)

// progressRe matches the percentage in yt-dlp download progress lines, e.g.:
// [download]  45.3% of 123.45MiB at 2.34MiB/s ETA 00:45
var progressRe = regexp.MustCompile(`\b(\d+(?:\.\d+)?)%`)

// ParseProgress extracts the download percentage from a yt-dlp stdout line.
// Returns (pct, true) for progress lines, (0, false) for everything else.
func ParseProgress(line string) (int, bool) {
	if !strings.Contains(line, "[download]") {
		return 0, false
	}
	m := progressRe.FindStringSubmatch(line)
	if m == nil {
		return 0, false
	}
	pct, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return 0, false
	}
	v := int(pct)
	if v < 0 || v > 100 {
		return 0, false
	}
	return v, true
}
