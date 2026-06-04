package manifesttest_test

import (
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
	decoded, err := manifesttest.Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if decoded.Module != "github.com/ZoneCNH/testkitx" || len(decoded.Gates) != 1 || len(decoded.Evidence) != 1 {
		t.Fatalf("unexpected decoded manifest: %+v", decoded)
	}
}
