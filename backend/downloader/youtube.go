package downloader

import "fmt"

func DownloadVideo(url, targetDir, newName string) {
	fmt.Printf("downloading video %s to %s/%s\n", url, targetDir, newName)
}
