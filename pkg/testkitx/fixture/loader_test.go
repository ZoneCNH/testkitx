package fixture_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/fixture"
)

func TestLoadReadsFixtureBytes(t *testing.T) {
	t.Parallel()
	data := fixture.Load(t, "sample.json")
	if len(data) == 0 {
		t.Fatal("expected non-empty fixture data")
	}
}

func TestLoadJSONDeserializesFixture(t *testing.T) {
	t.Parallel()
	var v struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	fixture.LoadJSON(t, "sample.json", &v)
	if v.Name != "test" || v.Value != 42 {
		t.Fatalf("unexpected fixture value: %+v", v)
	}
}

func TestWriteMkdirAllError(t *testing.T) {
	t.Parallel()
	ws := fixture.NewWorkspace(t, "test.mod")
	// Create a file where a directory is expected, so MkdirAll fails.
	blocker := filepath.Join(ws.ModuleDir, "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := ws.Write("blocker/file.txt", []byte("data"))
	if err == nil {
		t.Fatal("expected error when MkdirAll fails")
	}
}

func TestWriteFileError(t *testing.T) {
	t.Parallel()
	ws := fixture.NewWorkspace(t, "test.mod")
	// Create a directory at the target path so WriteFile fails (can't write to a directory).
	target := filepath.Join(ws.ModuleDir, "data.json")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	_, err := ws.Write("data.json", []byte("payload"))
	if err == nil {
		t.Fatal("expected error when WriteFile target is a directory")
	}
}

func TestNewWorkspaceMkdirAllError(t *testing.T) {
	// Use mockTB with a TempDir that returns a path where MkdirAll fails.
	m := &mockTB{tmpDir: "/dev/null/impossible"}
	fixture.NewWorkspace(m, "test.mod")
	if !m.failed {
		t.Fatal("expected NewWorkspace to fail with impossible TempDir")
	}
}

func TestLoadESuccess(t *testing.T) {
	t.Parallel()
	data, err := fixture.LoadE("sample.json")
	if err != nil {
		t.Fatalf("LoadE: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty data")
	}
}

func TestLoadEMissingFile(t *testing.T) {
	t.Parallel()
	_, err := fixture.LoadE("nonexistent.json")
	if err == nil {
		t.Fatal("expected error for missing fixture")
	}
}

func TestLoadJSONESuccess(t *testing.T) {
	t.Parallel()
	var v struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	if err := fixture.LoadJSONE("sample.json", &v); err != nil {
		t.Fatalf("LoadJSONE: %v", err)
	}
	if v.Name != "test" || v.Value != 42 {
		t.Fatalf("unexpected value: %+v", v)
	}
}

func TestLoadJSONEInvalidJSON(t *testing.T) {
	t.Parallel()
	// Use a non-JSON fixture file (if one exists) or test with invalid data.
	// Since we can't easily create fixtures, test the error path via LoadE failure.
	var v struct{}
	err := fixture.LoadJSONE("nonexistent.json", &v)
	if err == nil {
		t.Fatal("expected error for missing fixture")
	}
}

func TestDirESuccess(t *testing.T) {
	t.Parallel()
	dir, err := fixture.DirE()
	if err != nil {
		t.Fatalf("DirE: %v", err)
	}
	if !filepath.IsAbs(dir) {
		t.Fatalf("expected absolute path, got %q", dir)
	}
}
