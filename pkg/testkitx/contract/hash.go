package contract

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// HashFile computes the SHA256 hex digest of the file at path.
func HashFile(path string) (string, error) {
	return FileSHA256(path)
}

// HashDir computes a single SHA256 hex digest over all files in dir.
// Files are processed in sorted path order so the result is deterministic.
// Symlinks and subdirectories are walked recursively.
func HashDir(dir string) (string, error) {
	h := sha256.New()
	var paths []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(dir, path)
		paths = append(paths, rel)
		return nil
	})
	if err != nil {
		return "", err
	}
	sort.Strings(paths)
	for _, rel := range paths {
		// Include the relative path so renames are detected.
		h.Write([]byte(rel))
		h.Write([]byte{0})
		abs := filepath.Join(dir, rel)
		data, readErr := os.ReadFile(abs)
		if readErr != nil {
			return "", readErr
		}
		h.Write(data)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// AssertFileHash verifies that the file at path has the expected SHA256 digest.
// The test fails if the hashes do not match.
func AssertFileHash(t *testing.T, path, expected string) {
	t.Helper()
	actual, err := HashFile(path)
	if err != nil {
		t.Fatalf("hash file %s: %v", path, err)
	}
	if !strings.EqualFold(actual, expected) {
		t.Fatalf("hash mismatch for %s: got %s want %s", path, actual, expected)
	}
}
