package manifesttest_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/manifesttest"
)

func TestManifestRoundTripAndValidate(t *testing.T) {
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
	manifesttest.AssertChecksum(t, path)

	decoded, err := manifesttest.Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if decoded.Module != "github.com/ZoneCNH/testkitx" || len(decoded.Gates) != 1 || len(decoded.Evidence) != 1 {
		t.Fatalf("unexpected decoded manifest: %+v", decoded)
	}
}

func TestManifestChecksumDetectsDrift(t *testing.T) {
	manifest := manifesttest.New("github.com/ZoneCNH/testkitx", "abc123")
	path := filepath.Join(t.TempDir(), "manifest.json")
	if err := manifesttest.Write(path, manifest); err != nil {
		t.Fatal(err)
	}
	manifesttest.AssertChecksum(t, path)

	if err := os.WriteFile(path, []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := manifesttest.VerifyChecksum(path); err == nil {
		t.Fatal("VerifyChecksum() succeeded after manifest drift, want error")
	}
}
