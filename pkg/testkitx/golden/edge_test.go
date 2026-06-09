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
	// Create a directory at the target path so WriteFile fails.
	dir := filepath.Join(t.TempDir(), "golden.golden")
	os.MkdirAll(dir, 0o755)
	golden.AssertBytes(m, dir, []byte("data"))
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


func TestUpdateMkdirAllDirectError(t *testing.T) {
	t.Setenv(golden.UpdateEnv, "1")
	blocker := filepath.Join(t.TempDir(), "blocker")
	os.WriteFile(blocker, []byte("x"), 0o644)
	m := &mockTB{}
	golden.AssertBytes(m, filepath.Join(blocker, "sub", "file.golden"), []byte("data"))
	if !m.failed {
		t.Fatal("expected failure on MkdirAll error")
	}
}


func TestWriteGoldenMkdirAllError(t *testing.T) {
	t.Setenv(golden.UpdateEnv, "1")
	err := golden.WriteGolden("/dev/null/impossible/file.golden", []byte("data"))
	if err == nil {
		t.Fatal("expected MkdirAll error")
	}
}

func TestWriteGoldenWriteFileError(t *testing.T) {
	t.Setenv(golden.UpdateEnv, "1")
	// Create a directory at the target path so WriteFile fails (can't write to a directory).
	dir := filepath.Join(t.TempDir(), "file.golden")
	os.MkdirAll(dir, 0o755)
	err := golden.WriteGolden(dir, []byte("data"))
	if err == nil {
		t.Fatal("expected WriteFile error")
	}
}

func TestWriteGoldenDisabled(t *testing.T) {
	t.Parallel()
	// When env is not set, WriteGolden should return nil (no-op).
	err := golden.WriteGolden("/tmp/test.golden", []byte("data"))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestCheckBytesMkdirAllError(t *testing.T) {
	t.Setenv(golden.UpdateEnv, "1")
	_, err := golden.CheckBytes("/dev/null/impossible/file.golden", []byte("data"))
	if err == nil {
		t.Fatal("expected MkdirAll error")
	}
}

func TestCheckBytesWriteFileError(t *testing.T) {
	t.Setenv(golden.UpdateEnv, "1")
	dir := filepath.Join(t.TempDir(), "file.golden")
	os.MkdirAll(dir, 0o755)
	_, err := golden.CheckBytes(dir, []byte("data"))
	if err == nil {
		t.Fatal("expected WriteFile error")
	}
}

func TestCheckJSONMarshalError(t *testing.T) {
	t.Parallel()
	_, err := golden.CheckJSON(filepath.Join(t.TempDir(), "out.json"), func() {})
	if err == nil {
		t.Fatal("expected marshal error")
	}
}
