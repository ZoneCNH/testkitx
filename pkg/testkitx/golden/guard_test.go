package golden_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/golden"
)

func TestUpdateEnabled(t *testing.T) {
	t.Setenv(golden.UpdateEnv, "")
	if golden.UpdateEnabled() {
		t.Fatal("expected UpdateEnabled false when env is unset")
	}
	t.Setenv(golden.UpdateEnv, "1")
	if !golden.UpdateEnabled() {
		t.Fatal("expected UpdateEnabled true when env=1")
	}
}

func TestUpdateWritesOnlyWhenEnabled(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.golden")

	// Without UPDATE_GOLDEN — should be no-op.
	t.Setenv(golden.UpdateEnv, "")
	golden.Update(t, path, []byte("hello"))
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("file should not exist when UPDATE_GOLDEN is off")
	}

	// With UPDATE_GOLDEN — should write.
	t.Setenv(golden.UpdateEnv, "1")
	golden.Update(t, path, []byte("hello"))
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Fatalf("unexpected content: %q", data)
	}
}

func TestAssertDelegatesCorrectly(t *testing.T) {
	path := filepath.Join(t.TempDir(), "assert.golden")

	// Write via update mode.
	t.Setenv(golden.UpdateEnv, "1")
	golden.Assert(t, path, []byte("world"))

	// Compare via normal mode.
	t.Setenv(golden.UpdateEnv, "")
	golden.Assert(t, path, []byte("world"))
}
