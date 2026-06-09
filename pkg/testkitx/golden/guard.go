package golden

import (
	"fmt"
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
	if err := WriteGolden(path, actual); err != nil {
		t.Fatalf("%v", err)
	}
}

// WriteGolden writes actual to the golden file at path when
// TESTKITX_UPDATE_GOLDEN=1 is set. Returns an error instead of failing.
func WriteGolden(path string, actual []byte) error {
	if !UpdateEnabled() {
		return nil
	}
	clean := filepath.Clean(path)
	if err := os.MkdirAll(filepath.Dir(clean), 0o755); err != nil {
		return fmt.Errorf("create golden dir: %w", err)
	}
	if err := os.WriteFile(clean, actual, 0o644); err != nil {
		return fmt.Errorf("update golden %s: %w", path, err)
	}
	return nil
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
