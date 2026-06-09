package fixture_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/fixture"
)

func TestNewWorkspaceCreatesIsolatedHomeAndModule(t *testing.T) {
	t.Parallel()
	workspace := fixture.NewWorkspace(t, "example.test/module")
	if workspace.Root == "" || workspace.Home == "" || workspace.ModuleDir == "" {
		t.Fatalf("workspace paths not populated: %+v", workspace)
	}
	if workspace.Home == os.Getenv("HOME") {
		t.Fatalf("workspace reused process HOME: %s", workspace.Home)
	}
	if got := workspace.Env["HOME"]; got != workspace.Home {
		t.Fatalf("HOME env = %q, want %q", got, workspace.Home)
	}
	goMod, err := os.ReadFile(filepath.Join(workspace.ModuleDir, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(goMod), "module example.test/module") {
		t.Fatalf("unexpected go.mod: %s", goMod)
	}

	workspace.WriteOrFatal(t, "nested/file.txt", []byte("contents"))
	got, err := os.ReadFile(filepath.Join(workspace.ModuleDir, "nested/file.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "contents" {
		t.Fatalf("unexpected fixture file: %q", got)
	}
}

func TestWriteCreatesNestedDirectories(t *testing.T) {
	t.Parallel()
	ws := fixture.NewWorkspace(t, "example.test/module")
	path, err := ws.Write("deep/nested/dir/file.txt", []byte("deep content"))
	if err != nil {
		t.Fatalf("Write: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(data) != "deep content" {
		t.Fatalf("unexpected content: %q", data)
	}
}

func TestWriteOrFatalCallsFatalOnError(t *testing.T) {
	ws := fixture.NewWorkspace(t, "example.test/module")
	path := ws.WriteOrFatal(&mockTB{}, "ok.txt", []byte("ok"))
	if path == "" {
		t.Fatal("expected non-empty path")
	}
}

func TestNewWorkspaceWithModulePath(t *testing.T) {
	t.Parallel()
	ws := fixture.NewWorkspace(t, "my.test/module")
	if ws.Env["GOWORK"] != "off" {
		t.Fatalf("expected GOWORK=off, got %q", ws.Env["GOWORK"])
	}
	goMod, err := os.ReadFile(filepath.Join(ws.ModuleDir, "go.mod"))
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}
	if !strings.Contains(string(goMod), "module my.test/module") {
		t.Fatalf("unexpected go.mod: %s", goMod)
	}
}

