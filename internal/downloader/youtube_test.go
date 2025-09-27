package downloader_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/mcrors/ytd/internal/downloader"
)

func fakeCmd(ctx context.Context, name string, args ...string) *exec.Cmd {
	// Re-invoke the test binary in helper mode.
	// The `--` ensures everything after it is NOT parsed as go test flags.
	cmd := exec.CommandContext(ctx, os.Args[0],
		append([]string{"-test.run=TestHelperProcess", "--"}, args...)...)

	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	// Start with all args, then skip until the `--` separator.
	args := os.Args[1:]
	for len(args) > 0 && args[0] != "--" {
		args = args[1:]
	}
	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}

	// Now `args` contains exactly what your code passed: e.g. ["--no-warning","--print","%(channel)s",url]
	for i, a := range args {
		if a == "--print" && i+1 < len(args) && args[i+1] == "%(channel)s" {
			fmt.Fprintln(os.Stdout, "Flo Woelki")
			os.Exit(0)
		}
	}

	fmt.Fprintln(os.Stderr, "unexpected arguments:", args)
	os.Exit(2)
}

func TestGetChannel_ParsesOutput(t *testing.T) {
	// Given a YouTube
	yt := &downloader.YouTube{
		Bin:      "test-yt-dlp",
		Command:  fakeCmd,
		LookPath: func(string) (string, error) { return "/tmp/yt-dlp", nil },
	}
	got, err := yt.GetChannel(context.Background(), "https://x")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if got != "Flo Woelki" {
		t.Fatalf("got %q", got)
	}
}

func TestGetChannel_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network-dependent integration test in -short")
	}

	yt := downloader.NewYouTube()
	url := "https://www.youtube.com/watch?v=c8H0w4yBL10"
	want := "Flo Woelki"

	got, err := yt.GetChannel(context.Background(), url)

	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if got == "" {
		t.Fatalf("empty channel name")
	}
	if want != got {
		t.Errorf("wanted: %s, got: %s", want, got)
	}
}
