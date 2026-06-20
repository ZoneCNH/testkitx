package testkit

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// mockTB implements testing.TB for Go 1.26 without calling runtime.Goexit on Fatalf.
type mockTB struct {
	testing.TB
	failed bool
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
func (m *mockTB) TempDir() string                   { return os.TempDir() }
func (m *mockTB) Setenv(key, value string)          {}
func (m *mockTB) Cleanup(func())                    {}
func (m *mockTB) Error(args ...any)                 { m.failed = true }
func (m *mockTB) Fatal(args ...any)                 { m.failed = true }
func (m *mockTB) Fail()                             { m.failed = true }
func (m *mockTB) ArtifactDir() string               { return os.TempDir() }
func (m *mockTB) Attr(key, value string)            {}
func (m *mockTB) Chdir(dir string)                  {}
func (m *mockTB) Context() context.Context          { return context.Background() }
func (m *mockTB) Output() io.Writer                 { return io.Discard }

func TestRequireNoErrorNonNil(t *testing.T) {
	m := &mockTB{}
	RequireNoError(m, os.ErrNotExist)
	if !m.failed {
		t.Fatal("expected failure on non-nil error")
	}
}

func TestRequireGoldenMismatch(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "golden.golden")
	if err := os.WriteFile(path, []byte("expected"), 0o644); err != nil {
		t.Fatal(err)
	}
	m := &mockTB{}
	RequireGolden(m, path, []byte("actual"))
	if !m.failed {
		t.Fatal("expected failure on mismatch")
	}
}

func TestRequireGoldenReadError(t *testing.T) {
	m := &mockTB{}
	RequireGolden(m, "/nonexistent/golden.golden", []byte("data"))
	if !m.failed {
		t.Fatal("expected failure on read error")
	}
}
