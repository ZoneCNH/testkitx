package testkit

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRequireGoldenAcceptsMatchingContent(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.golden")

	if err := os.WriteFile(path, []byte("ok\n"), 0o600); err != nil {
		t.Fatalf("write golden: %v", err)
	}

	RequireGolden(t, path, []byte("ok\n"))
}
