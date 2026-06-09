package fixture

import (
	"encoding/json"
	"fmt"
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
	data, err := LoadE(path)
	if err != nil {
		t.Fatalf("load fixture: %v", err)
	}
	return data
}

// LoadE reads a fixture file, resolving the path relative to the caller's
// file. Returns an error instead of failing the test.
func LoadE(path string) ([]byte, error) {
	full, err := resolvePathE(2, path)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(full)
	if err != nil {
		return nil, fmt.Errorf("load fixture %s: %w", full, err)
	}
	return data, nil
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

// LoadJSONE reads and unmarshals a fixture file. Returns an error instead
// of failing the test.
func LoadJSONE(path string, v any) error {
	data, err := LoadE(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("unmarshal fixture %s: %w", path, err)
	}
	return nil
}

// Dir returns the absolute path to the testdata/fixtures directory relative
// to the calling test file's package. This is useful when you need to build
// paths programmatically.
func Dir(t *testing.T) string {
	t.Helper()
	dir, err := DirE()
	if err != nil {
		t.Fatalf("fixture.Dir: %v", err)
	}
	return dir
}

// DirE returns the testdata/fixtures directory path. Returns an error
// instead of failing the test.
func DirE() (string, error) {
	_, filename, _, ok := runtime.Caller(2)
	if !ok {
		return "", fmt.Errorf("unable to determine caller file")
	}
	return filepath.Join(filepath.Dir(filename), "testdata", "fixtures"), nil
}

// resolvePathE finds the absolute path for a fixture relative to the
// test file at the given caller skip depth.
func resolvePathE(skip int, path string) (string, error) {
	_, filename, _, ok := runtime.Caller(skip)
	if !ok {
		return "", fmt.Errorf("unable to determine caller file")
	}
	return filepath.Join(filepath.Dir(filename), "testdata", "fixtures", filepath.Clean(path)), nil
}
