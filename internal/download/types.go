package download

type DownloadCommand struct {
	TargetDir string
	URL       string
	NewName   string
}

type DownloadResult struct {
	Filename string
	Message  string
}

type DownloadError error
