package download_test

import (
	"testing"

	"github.com/mcrors/ytd/internal/download"
)

func TestFormatArgs(t *testing.T) {
	tests := []struct {
		format  download.Format
		wantErr bool
		wantArg string // spot-check the -f value
	}{
		{download.FormatBest, false, "bestvideo+bestaudio/best"},
		{download.Format1080p, false, "bestvideo[height<=1080]+bestaudio/best[height<=1080]"},
		{download.FormatAudio, false, "bestaudio/best"},
		{"invalid", true, ""},
		{"", true, ""},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			args, err := download.FormatArgs(tt.format)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FormatArgs(%q) error = %v, wantErr %v", tt.format, err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if len(args) < 2 || args[0] != "-f" {
				t.Fatalf("FormatArgs(%q) = %v, want first args to be [-f <format>]", tt.format, args)
			}
			if args[1] != tt.wantArg {
				t.Errorf("FormatArgs(%q) format arg = %q, want %q", tt.format, args[1], tt.wantArg)
			}
		})
	}
}

func TestFormatArgs_AudioExtractsMP3(t *testing.T) {
	args, err := download.FormatArgs(download.FormatAudio)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	has := func(flag string) bool {
		for _, a := range args {
			if a == flag {
				return true
			}
		}
		return false
	}
	if !has("--extract-audio") {
		t.Error("audio format missing --extract-audio flag")
	}
	if !has("mp3") {
		t.Error("audio format missing mp3 value")
	}
}
