package repotest_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/repotest"
)

func TestWriteFileCreatesParentDirectories(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	if _, err := repotest.WriteFile(root, "a/b/c.txt", []byte("data")); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(root, "a/b/c.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "data" {
		t.Fatalf("unexpected file content: %q", got)
	}
}
