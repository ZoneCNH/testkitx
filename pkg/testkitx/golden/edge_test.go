package golden_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/golden"
)

func TestCanonicalJSONMarshalError(t *testing.T) {
	t.Parallel()
	m := &mockTB{}
	golden.AssertJSON(m, filepath.Join(t.TempDir(), "out.json"), func() {})
	if !m.failed {
		t.Fatal("expected failure for non-marshalable value")
	}
}

func TestCanonicalJSONSuccess(t *testing.T) {
	path := filepath.Join(t.TempDir(), "golden.golden")
	t.Setenv(golden.UpdateEnv, "1")
	evidence := golden.AssertJSON(t, path, map[string]int{"a": 1, "b": 2})
	if !evidence.Matched {
		t.Fatal("expected Matched=true")
	}
}

func TestAssertBytesUpdatePath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "golden.golden")
	t.Setenv(golden.UpdateEnv, "1") //nolint:usetesting // cannot combine Setenv with Parallel
	evidence := golden.AssertBytes(t, path, []byte("updated content"))
	if !evidence.Updated {
		t.Fatal("expected Updated=true when env is set")
	}
}

func TestAssertBytesMatch(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "golden.golden")
	content := []byte("expected")
	os.MkdirAll(filepath.Dir(path), 0o755)
	os.WriteFile(path, content, 0o644)
	evidence := golden.AssertBytes(t, path, content)
	if !evidence.Matched {
		t.Fatal("expected Matched=true")
	}
}

func TestUpdateMkdirAllError(t *testing.T) {
	m := &mockTB{}
	t.Setenv(golden.UpdateEnv, "1") //nolint:usetesting // cannot combine Setenv with Parallel
	// Use AssertBytes which takes testing.TB - this covers the MkdirAll path
	golden.AssertBytes(m, "/nonexistent/dir/deep/file.golden", []byte("data"))
	if !m.failed {
		t.Fatal("expected failure on MkdirAll error")
	}
}

func TestAssertBytesUpdateWriteError(t *testing.T) {
	m := &mockTB{}
	t.Setenv(golden.UpdateEnv, "1") //nolint:usetesting // cannot combine Setenv with Parallel
	// Path where dir creation succeeds but write fails (read-only)
	// Use a path with a file where a directory is expected
	blocker := filepath.Join(t.TempDir(), "blocker")
	os.WriteFile(blocker, []byte("x"), 0o644)
	path := filepath.Join(blocker, "sub", "golden.golden")
	golden.AssertBytes(m, path, []byte("data"))
	if !m.failed {
		t.Fatal("expected failure on WriteFile error")
	}
}

func TestAssertJSONSuccessUpdate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "golden.golden")
	t.Setenv(golden.UpdateEnv, "1") //nolint:usetesting // cannot combine Setenv with Parallel
	evidence := golden.AssertJSON(t, path, map[string]string{"key": "value"})
	if !evidence.Matched {
		t.Fatal("expected Matched=true for JSON golden update")
	}
}
