package downloader_test

import (
	"context"
	"testing"

	"github.com/mcrors/ytd/internal/downloader"
)

func TestGetChannel(t *testing.T) {
	underTest := downloader.NewYouTube()
	url := "https://www.youtube.com/watch?v=c8H0w4yBL10"
	want := "Flo Woelki"

	got, err := underTest.GetChannel(context.Background(), url)

	if err != nil {
		t.Errorf("GetChannel error: %v", err)
	}

	if want != got {
		t.Errorf("wanted: %s, got: %s", want, got)
	}
}
