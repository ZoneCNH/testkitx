package golden

import (
	"os"
	"path/filepath"
	"testing"
)

// UpdateEnabled reports whether TESTKITX_UPDATE_GOLDEN=1 is set.
func UpdateEnabled() bool {
	return os.Getenv(UpdateEnv) == "1"
}

// Update writes actual to path when TESTKITX_UPDATE_GOLDEN=1 is set.
// It is a no-op otherwise.
func Update(t *testing.T, path string, actual []byte) {
	t.Helper()
	if !UpdateEnabled() {
		return
	}
	clean := filepath.Clean(path)
	if err := os.MkdirAll(filepath.Dir(clean), 0o755); err != nil {
		t.Fatalf("create golden dir: %v", err)
	}
	if err := os.WriteFile(clean, actual, 0o644); err != nil {
		t.Fatalf("update golden %s: %v", path, err)
	}
}

// Assert compares actual against the golden file at path. When
// TESTKITX_UPDATE_GOLDEN=1 is set it writes actual to path instead.
func Assert(t *testing.T, path string, actual []byte) {
	t.Helper()
	if UpdateEnabled() {
		Update(t, path, actual)
		return
	}
	AssertBytes(t, path, actual)
}
