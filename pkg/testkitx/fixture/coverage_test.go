package fixture_test

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/fixture"
)

// mockTB implements testing.TB for Go 1.26 without calling runtime.Goexit on Fatalf.
type mockTB struct {
	testing.TB
	failed bool
	tmpDir string
}

func (m *mockTB) Helper()                           {}
func (m *mockTB) Fatalf(format string, args ...any) { m.failed = true }
func (m *mockTB) Errorf(format string, args ...any) { m.failed = true }
func (m *mockTB) FailNow()                          { m.failed = true }
func (m *mockTB) Failed() bool                      { return m.failed }
func (m *mockTB) Name() string                      { return "mock" }
func (m *mockTB) Log(args ...any)                   {}
func (m *mockTB) Logf(format string, args ...any)   {}
func (m *mockTB) Skip(args ...any)                  {}
func (m *mockTB) Skipf(format string, args ...any)  {}
func (m *mockTB) SkipNow()                          {}
func (m *mockTB) Skipped() bool                     { return false }
func (m *mockTB) TempDir() string {
	if m.tmpDir != "" {
		return m.tmpDir
	}
	return os.TempDir()
}
func (m *mockTB) Setenv(key, value string) {}
func (m *mockTB) Cleanup(func())           {}
func (m *mockTB) Error(args ...any)        { m.failed = true }
func (m *mockTB) Fatal(args ...any)        { m.failed = true }
func (m *mockTB) Fail()                    { m.failed = true }
func (m *mockTB) ArtifactDir() string      { return os.TempDir() }
func (m *mockTB) Attr(key, value string)   {}
func (m *mockTB) Chdir(dir string)         {}
func (m *mockTB) Context() context.Context { return context.Background() }
func (m *mockTB) Output() io.Writer        { return io.Discard }

func TestNewWorkspaceEmptyModulePath(t *testing.T) {
	t.Parallel()
	ws := fixture.NewWorkspace(t, "")
	if ws.Root == "" {
		t.Fatal("expected non-empty root")
	}
	if _, err := os.Stat(filepath.Join(ws.ModuleDir, "go.mod")); !os.IsNotExist(err) {
		t.Fatal("go.mod should not exist when modulePath is empty")
	}
}

func TestWriteReturnsPath(t *testing.T) {
	t.Parallel()
	ws := fixture.NewWorkspace(t, "example.mod")
	path, err := ws.Write("sub/file.txt", []byte("content"))
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "content" {
		t.Fatalf("unexpected content: %q", data)
	}
}

func TestDir(t *testing.T) {
	t.Parallel()
	dir := fixture.Dir(t)
	if !filepath.IsAbs(dir) {
		t.Fatalf("expected absolute path, got %q", dir)
	}
	if filepath.Base(dir) != "fixtures" {
		t.Fatalf("expected basename 'fixtures', got %q", filepath.Base(dir))
	}
}
