package download

import (
	"context"
	"path/filepath"
)

type Downloader interface {
	Download(ctx context.Context, url, targetDir, newName string) error
	GetChannel(ctx context.Context, url string) (string, error)
}

type downloadService struct {
	baseDir    string
	downloader Downloader
}

func NewDownloadService(baseDir string, downloader Downloader) *downloadService {
	return &downloadService{
		baseDir:    baseDir,
		downloader: downloader,
	}
}

func (ds *downloadService) Download(ctx context.Context, dc DownloadCommand) (*DownloadResult, error) {
	rel, err := normalizeTwoLevel(dc.TargetDir)
	if err != nil {
		return nil, err
	}

	channelName, err := ds.downloader.GetChannel(ctx, dc.URL)
	if err != nil {
		return nil, err
	}

	target := filepath.Join(ds.baseDir, rel, channelName)

	if err := ds.downloader.Download(ctx, dc.URL, target, dc.NewName); err != nil {
		return nil, err
	}

	filename := filepath.Join(target, dc.NewName)
	return &DownloadResult{
		Filename: filename,
		Message:  "Download Completed Succesfully",
	}, nil
}
