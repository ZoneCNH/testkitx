package fixture

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// Load reads the fixture file at path and returns its bytes.
// The path is relative to the testdata/fixtures directory next to the calling
// test file. If the file does not exist the test fails immediately.
func Load(t *testing.T, path string) []byte {
	t.Helper()
	full := resolvePath(t, path)
	data, err := os.ReadFile(full)
	if err != nil {
		t.Fatalf("load fixture %s: %v", full, err)
	}
	return data
}

// LoadJSON reads the fixture file at path and unmarshals it into v.
// The path is relative to the testdata/fixtures directory next to the calling
// test file.
func LoadJSON(t *testing.T, path string, v any) {
	t.Helper()
	data := Load(t, path)
	if err := json.Unmarshal(data, v); err != nil {
		t.Fatalf("unmarshal fixture %s: %v", path, err)
	}
}

// Dir returns the absolute path to the testdata/fixtures directory relative
// to the calling test file's package. This is useful when you need to build
// paths programmatically.
func Dir(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		t.Fatal("fixture.Dir: unable to determine caller file")
	}
	return filepath.Join(filepath.Dir(filename), "testdata", "fixtures")
}

// resolvePath finds the absolute path for a fixture relative to the caller's
// test file location.
func resolvePath(t *testing.T, path string) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(2)
	if !ok {
		t.Fatal("fixture: unable to determine caller file")
	}
	return filepath.Join(filepath.Dir(filename), "testdata", "fixtures", filepath.Clean(path))
}
