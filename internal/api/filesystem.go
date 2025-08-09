package api

import "os"

func findDirs(entries []os.DirEntry) []string {
	// TODO: should this be recursive, so we can sub-dirs
	var results []string
	for _, entry := range entries {
		if entry.IsDir() {
			results = append(results, entry.Name())
		}
	}
	return results
}
