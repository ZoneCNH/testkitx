package repotest_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/repotest"
)

func TestWriteFileReturnsCorrectPath(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	path, err := repotest.WriteFile(root, "file.txt", []byte("hello"))
	if err != nil {
		t.Fatal(err)
	}
	expected := filepath.Join(root, "file.txt")
	if path != expected {
		t.Fatalf("expected path %q, got %q", expected, path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Fatalf("unexpected content: %q", data)
	}
}
