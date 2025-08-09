package api

type DownloadRequest struct {
	URL       string `json:"url"`
	TargetDir string `json:"targetDir"`
	NewName   string `json:"newName"`
}

type CreateDirectoryRequest struct {
	Dir string `json:"dir"`
}

type DirectoriesResponse struct {
	Directories []string `json:"directories"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
