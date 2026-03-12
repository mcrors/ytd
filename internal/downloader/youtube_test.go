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

	output, ok := os.LookupEnv("FAKE_YTDLP_STDOUT")
	if !ok {
		fmt.Println("error: couldn't find FAKE_YTDLP_STDOUT in env")
	}

	// Now `args` contains exactly what your code passed: e.g. ["--no-warning","--print","%(channel)s",url]
	for i, a := range args {
		if a == "--print" && i+1 < len(args) && args[i+1] == "%(channel)s" {
			fmt.Fprintln(os.Stdout, output)
			os.Exit(0)
		}
	}

	fmt.Fprintln(os.Stderr, "unexpected arguments:", args)
	os.Exit(2)
}

var mockedYtDlp = downloader.NewYouTube(
	"test-yt-dlp",
	fakeCmd,
	func(file string) (string, error) { return "/tmp/yt-dlp", nil },
)

var getChannelTests = []struct {
	name   string
	output string
	want   string
}{
	{"normal", "Flo Woelki", "Flo Woelki"},
	{"spaces", "  Flo Woelki  ", "Flo Woelki"},
	{"new lines", "Flo Woelki\n\n", "Flo Woelki"},
}

func TestGetChannel_ParsesOutput(t *testing.T) {
	for _, tt := range getChannelTests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("FAKE_YTDLP_STDOUT", tt.output)
			got, err := mockedYtDlp.GetChannel(context.Background(), "https://x")
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			if got != tt.want {
				t.Errorf("wanted: %s, got: %s", tt.want, got)
			}
		})
	}
}

func TestGetChannel_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network-dependent integration test in -short")
	}

	yt := downloader.NewYouTube("yt-dlp", exec.CommandContext, exec.LookPath)
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
