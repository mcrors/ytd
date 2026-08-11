package download

type DownloadCommand struct {
	TargetDir string
	URL       string
	NewName   string
	Format    Format
}

type DownloadResult struct {
	Filename string
	Message  string
}

type DownloadError error
