package fixture_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/fixture"
)

func TestNewWorkspaceCreatesIsolatedHomeAndModule(t *testing.T) {
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

	workspace.Write("nested/file.txt", []byte("contents"))
	got, err := os.ReadFile(filepath.Join(workspace.ModuleDir, "nested/file.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "contents" {
		t.Fatalf("unexpected fixture file: %q", got)
	}
}
