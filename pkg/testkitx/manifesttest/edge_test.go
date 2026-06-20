package manifesttest_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/manifesttest"
)

func TestWriteMkdirAllError(t *testing.T) {
	t.Parallel()
	blocker := filepath.Join(t.TempDir(), "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(blocker, "sub", "manifest.json")
	manifest := manifesttest.New("mod", "abc123")
	err := manifesttest.Write(path, manifest)
	if err == nil {
		t.Fatal("expected MkdirAll error when parent is a file")
	}
}

func TestVerifyChecksumReferenceMismatch(t *testing.T) {
	t.Parallel()
	manifest := manifesttest.New("mod", "abc123")
	path := filepath.Join(t.TempDir(), "manifest.json")
	if err := manifesttest.Write(path, manifest); err != nil {
		t.Fatal(err)
	}
	checksumPath := manifesttest.ChecksumPath(path)
	if err := os.WriteFile(checksumPath, []byte("abc123  wrong_name.json\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := manifesttest.VerifyChecksum(path, checksumPath)
	if err == nil {
		t.Fatal("expected error for reference mismatch")
	}
}

func TestVerifyChecksumNonExistentChecksum(t *testing.T) {
	t.Parallel()
	manifest := manifesttest.New("mod", "abc123")
	path := filepath.Join(t.TempDir(), "manifest.json")
	if err := manifesttest.Write(path, manifest); err != nil {
		t.Fatal(err)
	}
	err := manifesttest.VerifyChecksum(path, "/nonexistent/checksum.sha256")
	if err == nil {
		t.Fatal("expected error for nonexistent checksum file")
	}
}

func TestVerifyChecksumNonExistentManifest(t *testing.T) {
	t.Parallel()
	err := manifesttest.VerifyChecksum("/nonexistent/manifest.json", "/nonexistent/checksum.sha256")
	if err == nil {
		t.Fatal("expected error for nonexistent manifest file")
	}
}

func TestWriteChecksumNonExistentManifest(t *testing.T) {
	t.Parallel()
	err := manifesttest.WriteChecksum("/nonexistent/manifest.json", filepath.Join(t.TempDir(), "out.sha256"))
	if err == nil {
		t.Fatal("expected error for nonexistent manifest")
	}
}

func TestWriteInvalidManifest(t *testing.T) {
	t.Parallel()
	invalid := manifesttest.Manifest{}
	err := manifesttest.Write(filepath.Join(t.TempDir(), "m.json"), invalid)
	if err == nil || err.Error() != "manifest missing required fields" {
		t.Fatalf("expected validate error, got %v", err)
	}
}

func TestReadNonExistentManifest(t *testing.T) {
	t.Parallel()
	_, err := manifesttest.Read("/nonexistent/manifest.json")
	if err == nil {
		t.Fatal("expected error for nonexistent manifest")
	}
}
