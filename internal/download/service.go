package download

import (
	"context"
	"path/filepath"

	"github.com/mcrors/ytd/internal/pathutil"
)

type Downloader interface {
	Download(ctx context.Context, url, targetDir, newName string, format Format) error
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
	target, err := pathutil.SafeJoin(ds.baseDir, dc.TargetDir)
	if err != nil {
		return nil, err
	}

	if err := ds.downloader.Download(ctx, dc.URL, target, dc.NewName, dc.Format); err != nil {
		return nil, err
	}

	return &DownloadResult{
		Filename: filepath.Join(target, dc.NewName),
		Message:  "Download completed successfully",
	}, nil
}
