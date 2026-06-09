package fixture_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/fixture"
)

func TestNewWorkspaceEmptyModule(t *testing.T) {
	t.Parallel()
	ws := fixture.NewWorkspace(t, "")
	if ws.Root == "" || ws.Home == "" || ws.ModuleDir == "" {
		t.Fatal("expected non-empty workspace paths")
	}
}

func TestWriteOrFatalCallsFatal(t *testing.T) {
	t.Parallel()
	ws := fixture.NewWorkspace(t, "testmod")
	m := &mockTB{}
	blocker := filepath.Join(ws.ModuleDir, "block")
	os.WriteFile(blocker, []byte("x"), 0o644)
	ws.WriteOrFatal(m, "block/deep/file.txt", []byte("data"))
	if !m.failed {
		t.Fatal("expected WriteOrFatal to fail on MkdirAll error")
	}
}

func TestWriteOrFatalSuccess(t *testing.T) {
	t.Parallel()
	ws := fixture.NewWorkspace(t, "testmod")
	path := ws.WriteOrFatal(t, "hello.txt", []byte("world"))
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(data) != "world" {
		t.Fatalf("unexpected content: %q", data)
	}
}

func TestDirReturnsNonEmpty(t *testing.T) {
	dir := fixture.Dir(t)
	if dir == "" {
		t.Fatal("expected non-empty Dir result")
	}
	if !filepath.IsAbs(dir) {
		t.Fatalf("expected absolute path, got %q", dir)
	}
}
