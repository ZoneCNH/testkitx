// Package fixture creates isolated test workspaces.
package fixture

import (
	"os"
	"path/filepath"
	"testing"
)

type Workspace struct {
	Root      string
	Home      string
	ModuleDir string
	Env       map[string]string
}

func NewWorkspace(t testing.TB, modulePath string) Workspace {
	t.Helper()
	root := t.TempDir()
	home := filepath.Join(root, "home")
	mod := filepath.Join(root, "module")
	for _, dir := range []string{home, mod} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("create fixture dir: %v", err)
		}
	}
	if modulePath != "" {
		if err := os.WriteFile(filepath.Join(mod, "go.mod"), []byte("module "+modulePath+"\n\ngo 1.24\n"), 0o644); err != nil {
			t.Fatalf("write go.mod: %v", err)
		}
	}
	return Workspace{Root: root, Home: home, ModuleDir: mod, Env: map[string]string{"HOME": home, "GOWORK": "off"}}
}

func (w Workspace) Write(path string, data []byte) (string, error) {
	clean := filepath.Join(w.ModuleDir, filepath.Clean(path))
	if err := os.MkdirAll(filepath.Dir(clean), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(clean, data, 0o644); err != nil {
		return "", err
	}
	return clean, nil
}

// WriteOrFatal writes data to the workspace, calling t.Fatal on error.
func (w Workspace) WriteOrFatal(t testing.TB, path string, data []byte) string {
	t.Helper()
	clean, err := w.Write(path, data)
	if err != nil {
		t.Fatalf("write fixture file %s: %v", path, err)
	}
	return clean
}
