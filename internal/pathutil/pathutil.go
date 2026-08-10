package pathutil

import (
	"errors"
	"path/filepath"
	"strings"
)

// SafeJoin joins baseDir and userPath, then verifies the result stays under
// baseDir. Returns an error for traversal attempts or absolute userPaths.
// An empty userPath returns baseDir.
func SafeJoin(baseDir, userPath string) (string, error) {
	if filepath.IsAbs(userPath) {
		return "", errors.New("path must be relative")
	}

	joined := filepath.Join(baseDir, userPath)

	rel, err := filepath.Rel(baseDir, joined)
	if err != nil {
		return "", errors.New("invalid path")
	}
	if strings.HasPrefix(rel, "..") {
		return "", errors.New("path escapes base directory")
	}

	return joined, nil
}
