package api

import (
	"encoding/json"
	"net/http"
	"os"
)

func respondJson(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJson(w, status, ErrorResponse{Error: message})
}

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
