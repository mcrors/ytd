package pathutil_test

import (
	"testing"

	"github.com/mcrors/ytd/internal/pathutil"
)

func TestSafeJoin(t *testing.T) {
	base := "/media"

	tests := []struct {
		name     string
		userPath string
		want     string
		wantErr  bool
	}{
		{"normal path", "history/mary-beard", "/media/history/mary-beard", false},
		{"traversal blocked", "../../etc/passwd", "", true},
		{"absolute path blocked", "/etc/passwd", "", true},
		{"empty path returns base", "", "/media", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pathutil.SafeJoin(base, tt.userPath)
			if (err != nil) != tt.wantErr {
				t.Fatalf("SafeJoin() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Errorf("SafeJoin() = %q, want %q", got, tt.want)
			}
		})
	}
}
