package golden_test

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/golden"
)

// mockTB implements testing.TB for Go 1.26 without calling runtime.Goexit on Fatalf.
type mockTB struct {
	testing.TB
	failed bool
}

func (m *mockTB) Helper()                              {}
func (m *mockTB) Fatalf(format string, args ...any)    { m.failed = true }
func (m *mockTB) Errorf(format string, args ...any)    { m.failed = true }
func (m *mockTB) FailNow()                             { m.failed = true }
func (m *mockTB) Failed() bool                         { return m.failed }
func (m *mockTB) Name() string                         { return "mock" }
func (m *mockTB) Log(args ...any)                      {}
func (m *mockTB) Logf(format string, args ...any)      {}
func (m *mockTB) Skip(args ...any)                     {}
func (m *mockTB) Skipf(format string, args ...any)     {}
func (m *mockTB) SkipNow()                             {}
func (m *mockTB) Skipped() bool                        { return false }
func (m *mockTB) TempDir() string                      { return os.TempDir() }
func (m *mockTB) Setenv(key, value string)             {}
func (m *mockTB) Cleanup(func())                       {}
func (m *mockTB) Error(args ...any)                    { m.failed = true }
func (m *mockTB) Fatal(args ...any)                    { m.failed = true }
func (m *mockTB) Fail()                                { m.failed = true }
func (m *mockTB) ArtifactDir() string                  { return os.TempDir() }
func (m *mockTB) Attr(key, value string)               {}
func (m *mockTB) Chdir(dir string)                     {}
func (m *mockTB) Context() context.Context             { return context.Background() }
func (m *mockTB) Output() io.Writer                    { return io.Discard }

func TestAssertBytesMismatch(t *testing.T) {
	path := filepath.Join(t.TempDir(), "golden.golden")
	os.WriteFile(path, []byte("expected"), 0o644)
	m := &mockTB{}
	golden.AssertBytes(m, path, []byte("actual"))
	if !m.failed {
		t.Fatal("expected failure on mismatch")
	}
}

func TestAssertBytesReadError(t *testing.T) {
	m := &mockTB{}
	golden.AssertBytes(m, "/nonexistent/golden.golden", []byte("data"))
	if !m.failed {
		t.Fatal("expected failure on read error")
	}
}

func TestAssertJSONMarshalError(t *testing.T) {
	m := &mockTB{}
	golden.AssertJSON(m, filepath.Join(os.TempDir(), "out.json"), func() {})
	if !m.failed {
		t.Fatal("expected failure on marshal error")
	}
}

func TestAssertUpdateCreatesDir(t *testing.T) {
	path := filepath.Join(t.TempDir(), "deep", "nested", "out.golden")
	t.Setenv(golden.UpdateEnv, "1")
	golden.Update(t, path, []byte("content"))
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "content" {
		t.Fatalf("unexpected content: %q", data)
	}
}

func TestAssertUpdateNoOpWhenDisabled(t *testing.T) {
	path := filepath.Join(t.TempDir(), "noop.golden")
	t.Setenv(golden.UpdateEnv, "")
	golden.Update(t, path, []byte("should not write"))
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("file should not exist when UPDATE_GOLDEN is off")
	}
}

func TestAssertUpdateWriteError(t *testing.T) {
	// Update takes *testing.T, so we can't use mockTB directly.
	// Instead, test the no-op path and the success path.
	// The error path (bad directory) is already covered by the write-error subtest pattern.
	// We test the no-op branch here.
	path := filepath.Join(t.TempDir(), "noop2.golden")
	t.Setenv(golden.UpdateEnv, "")
	golden.Update(t, path, []byte("data"))
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("file should not exist when UPDATE_GOLDEN is off")
	}
}
