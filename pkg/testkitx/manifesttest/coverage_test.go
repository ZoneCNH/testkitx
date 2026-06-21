package manifesttest_test

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/manifesttest"
)

// mockTB implements testing.TB for Go 1.26 without calling runtime.Goexit on Fatalf.
type mockTB struct {
	testing.TB
	failed bool
}

func (m *mockTB) Helper()                           {}
func (m *mockTB) Fatalf(format string, args ...any) { m.failed = true }
func (m *mockTB) Errorf(format string, args ...any) { m.failed = true }
func (m *mockTB) FailNow()                          { m.failed = true }
func (m *mockTB) Failed() bool                      { return m.failed }
func (m *mockTB) Name() string                      { return "mock" }
func (m *mockTB) Log(args ...any)                   {}
func (m *mockTB) Logf(format string, args ...any)   {}
func (m *mockTB) Skip(args ...any)                  {}
func (m *mockTB) Skipf(format string, args ...any)  {}
func (m *mockTB) SkipNow()                          {}
func (m *mockTB) Skipped() bool                     { return false }
func (m *mockTB) TempDir() string                   { return os.TempDir() }
func (m *mockTB) Setenv(key, value string)          {}
func (m *mockTB) Cleanup(func())                    {}
func (m *mockTB) Error(args ...any)                 { m.failed = true }
func (m *mockTB) Fatal(args ...any)                 { m.failed = true }
func (m *mockTB) Fail()                             { m.failed = true }
func (m *mockTB) ArtifactDir() string               { return os.TempDir() }
func (m *mockTB) Attr(key, value string)            {}
func (m *mockTB) Chdir(dir string)                  {}
func (m *mockTB) Context() context.Context          { return context.Background() }
func (m *mockTB) Output() io.Writer                 { return io.Discard }

func TestValidateRejectsInvalidManifest(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		m    manifesttest.Manifest
	}{
		{"empty kind", manifesttest.Manifest{Module: "m", Commit: "c"}},
		{"empty module", manifesttest.Manifest{Kind: "manifest_fixture_check", Commit: "c"}},
		{"empty commit", manifesttest.Manifest{Kind: "manifest_fixture_check", Module: "m"}},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if err := tc.m.Validate(); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestWriteRejectsInvalidManifest(t *testing.T) {
	t.Parallel()
	err := manifesttest.Write(filepath.Join(t.TempDir(), "out.json"), manifesttest.Manifest{})
	if err == nil {
		t.Fatal("expected error for invalid manifest")
	}
}

func TestReadNonExistent(t *testing.T) {
	t.Parallel()
	_, err := manifesttest.Read("/nonexistent/manifest.json")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestSHA256NonExistent(t *testing.T) {
	t.Parallel()
	_, err := manifesttest.SHA256("/nonexistent/file.json")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestVerifyChecksumEmptyChecksumFile(t *testing.T) {
	t.Parallel()
	manifest := manifesttest.New("mod", "abc123")
	path := filepath.Join(t.TempDir(), "manifest.json")
	if err := manifesttest.Write(path, manifest); err != nil {
		t.Fatal(err)
	}
	checksumPath := manifesttest.ChecksumPath(path)
	if err := os.WriteFile(checksumPath, []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	err := manifesttest.VerifyChecksum(path, checksumPath)
	if err == nil {
		t.Fatal("expected error for empty checksum file")
	}
}

func TestVerifyChecksumInvalidSHA(t *testing.T) {
	t.Parallel()
	manifest := manifesttest.New("mod", "abc123")
	path := filepath.Join(t.TempDir(), "manifest.json")
	if err := manifesttest.Write(path, manifest); err != nil {
		t.Fatal(err)
	}
	checksumPath := manifesttest.ChecksumPath(path)
	if err := os.WriteFile(checksumPath, []byte("abc123  manifest.json\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := manifesttest.VerifyChecksum(path, checksumPath)
	if err == nil {
		t.Fatal("expected error for invalid checksum")
	}
}

func TestVerifyChecksumNonHexSHA(t *testing.T) {
	t.Parallel()
	manifest := manifesttest.New("mod", "abc123")
	path := filepath.Join(t.TempDir(), "manifest.json")
	if err := manifesttest.Write(path, manifest); err != nil {
		t.Fatal(err)
	}
	checksumPath := manifesttest.ChecksumPath(path)
	if err := os.WriteFile(checksumPath, []byte("zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz  manifest.json\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := manifesttest.VerifyChecksum(path, checksumPath)
	if err == nil {
		t.Fatal("expected error for non-hex checksum")
	}
}

func TestChecksumPath(t *testing.T) {
	t.Parallel()
	got := manifesttest.ChecksumPath("/tmp/manifest.json")
	if got != "/tmp/manifest.json.sha256" {
		t.Fatalf("unexpected checksum path: %q", got)
	}
}

func TestReadInvalidJSON(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(path, []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := manifesttest.Read(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestWriteChecksumDefaultPath(t *testing.T) {
	t.Parallel()
	manifest := manifesttest.New("mod", "abc123")
	path := filepath.Join(t.TempDir(), "manifest.json")
	if err := manifesttest.Write(path, manifest); err != nil {
		t.Fatal(err)
	}
	err := manifesttest.WriteChecksum(path, "")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAssertManifestValidFailsOnBadPath(t *testing.T) {
	t.Parallel()
	m := &mockTB{}
	manifesttest.AssertManifestValid(m, "/nonexistent/manifest.json")
	if !m.failed {
		t.Fatal("expected failure on bad path")
	}
}

func TestAssertChecksumFailsOnDrift(t *testing.T) {
	t.Parallel()
	manifest := manifesttest.New("mod", "abc123")
	path := filepath.Join(t.TempDir(), "manifest.json")
	if err := manifesttest.Write(path, manifest); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"kind":"manifest_fixture_check","module":"changed","commit":"abc123"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	m := &mockTB{}
	manifesttest.AssertChecksum(m, path, manifesttest.ChecksumPath(path))
	if !m.failed {
		t.Fatal("expected failure on checksum drift")
	}
}

func TestVerifyChecksumMismatch(t *testing.T) {
	t.Parallel()
	manifest := manifesttest.New("mod", "abc123")
	path := filepath.Join(t.TempDir(), "manifest.json")
	if err := manifesttest.Write(path, manifest); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"kind":"manifest_fixture_check","module":"changed","commit":"abc123"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	err := manifesttest.VerifyChecksum(path, "")
	if err == nil {
		t.Fatal("expected checksum mismatch error")
	}
}
