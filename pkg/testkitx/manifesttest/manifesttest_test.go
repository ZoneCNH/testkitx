package manifesttest_test

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/manifesttest"
)

func TestManifestRoundTripAndValidate(t *testing.T) {
	t.Parallel()
	manifest := manifesttest.New("github.com/ZoneCNH/testkitx", "abc123")
	manifest.Gates["test"] = "go test ./..."
	manifest.Evidence = append(manifest.Evidence, "evidence.json")
	if err := manifest.Validate(); err != nil {
		t.Fatalf("Validate() failed: %v", err)
	}
	path := filepath.Join(t.TempDir(), "manifest.json")
	if err := manifesttest.Write(path, manifest); err != nil {
		t.Fatal(err)
	}
	manifesttest.AssertManifestValid(t, path)
	manifesttest.AssertChecksum(t, path, manifesttest.ChecksumPath(path))
	decoded, err := manifesttest.Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if decoded.Module != "github.com/ZoneCNH/testkitx" || len(decoded.Gates) != 1 || len(decoded.Evidence) != 1 {
		t.Fatalf("unexpected decoded manifest: %+v", decoded)
	}
}

func TestManifestChecksumRejectsWrongFilename(t *testing.T) {
	t.Parallel()
	manifest := manifesttest.New("github.com/ZoneCNH/testkitx", "abc123")
	path := filepath.Join(t.TempDir(), "manifest.json")
	if err := manifesttest.Write(path, manifest); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(data)
	checksumPath := manifesttest.ChecksumPath(path)
	if err := os.WriteFile(checksumPath, []byte(hex.EncodeToString(sum[:])+"  wrong.json\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	err = manifesttest.VerifyChecksum(path, checksumPath)
	if err == nil || !strings.Contains(err.Error(), `references "wrong.json", want "manifest.json"`) {
		t.Fatalf("expected filename mismatch, got %v", err)
	}
}

func TestManifestChecksumDetectsDrift(t *testing.T) {
	t.Parallel()
	manifest := manifesttest.New("github.com/ZoneCNH/testkitx", "abc123")
	path := filepath.Join(t.TempDir(), "manifest.json")
	if err := manifesttest.Write(path, manifest); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"kind":"manifest_fixture_check","module":"changed","commit":"abc123"}`+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := manifesttest.VerifyChecksum(path, ""); err == nil {
		t.Fatal("VerifyChecksum() succeeded after manifest drift")
	}
}

func TestReadNonExistentFile(t *testing.T) {
	t.Parallel()
	_, err := manifesttest.Read("/nonexistent/manifest.json")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestWriteChecksumSHA256Error(t *testing.T) {
	t.Parallel()
	// SHA256 on nonexistent file -> WriteChecksum should fail
	err := manifesttest.WriteChecksum("/nonexistent/manifest.json", "")
	if err == nil {
		t.Fatal("expected error for nonexistent manifest")
	}
}

func TestVerifyChecksumEmptyFile(t *testing.T) {
	t.Parallel()
	manifest := manifesttest.New("github.com/ZoneCNH/testkitx", "abc123")
	path := filepath.Join(t.TempDir(), "manifest.json")
	if err := manifesttest.Write(path, manifest); err != nil {
		t.Fatal(err)
	}
	checksumPath := manifesttest.ChecksumPath(path)
	if err := os.WriteFile(checksumPath, []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	err := manifesttest.VerifyChecksum(path, checksumPath)
	if err == nil || !strings.Contains(err.Error(), "empty checksum file") {
		t.Fatalf("expected empty checksum error, got %v", err)
	}
}

func TestVerifyChecksumInvalidHex(t *testing.T) {
	t.Parallel()
	manifest := manifesttest.New("github.com/ZoneCNH/testkitx", "abc123")
	path := filepath.Join(t.TempDir(), "manifest.json")
	if err := manifesttest.Write(path, manifest); err != nil {
		t.Fatal(err)
	}
	checksumPath := manifesttest.ChecksumPath(path)
	if err := os.WriteFile(checksumPath, []byte(strings.Repeat("z", 64)+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := manifesttest.VerifyChecksum(path, checksumPath)
	if err == nil || !strings.Contains(err.Error(), "invalid sha256") {
		t.Fatalf("expected invalid sha256 error, got %v", err)
	}
}

func TestVerifyChecksumNonExistentChecksumFile(t *testing.T) {
	t.Parallel()
	manifest := manifesttest.New("github.com/ZoneCNH/testkitx", "abc123")
	path := filepath.Join(t.TempDir(), "manifest.json")
	if err := manifesttest.Write(path, manifest); err != nil {
		t.Fatal(err)
	}
	err := manifesttest.VerifyChecksum(path, "/nonexistent/checksum.sha256")
	if err == nil {
		t.Fatal("expected error for nonexistent checksum file")
	}
}
