package golden_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/golden"
)

func TestAssertBytesUpdatesOnlyWhenOptedIn(t *testing.T) {
	path := filepath.Join(t.TempDir(), "actual.golden")

	t.Setenv(golden.UpdateEnv, "1")
	evidence := golden.AssertBytes(t, path, []byte("hello"))
	if !evidence.Updated || !evidence.Matched || evidence.Kind != "golden_check" {
		t.Fatalf("unexpected update evidence: %+v", evidence)
	}

	t.Setenv(golden.UpdateEnv, "")
	evidence = golden.AssertBytes(t, path, []byte("hello"))
	if evidence.Updated || !evidence.Matched || evidence.ActualSHA256 == "" {
		t.Fatalf("unexpected compare evidence: %+v", evidence)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello" {
		t.Fatalf("golden changed without opt-in: %q", got)
	}
}

func TestAssertJSONCanonicalizesOutput(t *testing.T) {
	path := filepath.Join(t.TempDir(), "value.json")
	t.Setenv(golden.UpdateEnv, "1")
	golden.AssertJSON(t, path, map[string]any{"b": 2, "a": 1})

	t.Setenv(golden.UpdateEnv, "")
	golden.AssertBytes(t, path, []byte("{\n  \"a\": 1,\n  \"b\": 2\n}"))
}

func TestAssertBytesCreatesDirectoryInUpdateMode(t *testing.T) {
	path := filepath.Join(t.TempDir(), "deep", "nested", "dir", "golden.json")
	t.Setenv(golden.UpdateEnv, "1")
	evidence := golden.AssertBytes(t, path, []byte(`{"ok":true}`))
	if !evidence.Updated {
		t.Fatal("expected Updated=true")
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{"ok":true}` {
		t.Fatalf("content = %q, want %q", got, `{"ok":true}`)
	}
}
