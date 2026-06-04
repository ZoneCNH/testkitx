// Package repotest contains repository fixture helpers.
package repotest

import (
	"os"
	"path/filepath"
)

func WriteFile(root, rel string, data []byte) (string, error) {
	path := filepath.Join(root, filepath.Clean(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	return path, os.WriteFile(path, data, 0o644)
}
