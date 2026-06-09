package golden_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/golden"
)

func TestCheckBytesMissingFile(t *testing.T) {
	_, err := golden.CheckBytes("/tmp/nonexistent-golden-12345.bin", []byte("data"))
	if err == nil {
		t.Fatal("expected error for missing golden file")
	}
}

func TestCheckJSONCanonicalError(t *testing.T) {
	ch := make(chan int)
	_, err := golden.CheckJSON("/tmp/nonexistent.json", ch)
	if err == nil {
		t.Fatal("expected error for unmarshalable value")
	}
}

func TestWriteGoldenCreatesDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "golden.bin")
	t.Setenv(golden.UpdateEnv, "1")
	err := golden.WriteGolden(path, []byte("hello"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	if string(data) != "hello" {
		t.Fatalf("got %q", data)
	}
}

func TestWriteGoldenNoopWhenDisabled(t *testing.T) {
	t.Setenv(golden.UpdateEnv, "")
	err := golden.WriteGolden("/tmp/should-not-exist-12345.bin", []byte("data"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWriteGoldenInvalidPath(t *testing.T) {
	t.Setenv(golden.UpdateEnv, "1")
	err := golden.WriteGolden("/nonexistent-root-dir-12345/golden.bin", []byte("data"))
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}
