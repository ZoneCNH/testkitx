package boundarytest_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/boundarytest"
)

func TestScanProductionImportsFlagsOnlyProductionViolations(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	write := func(rel, body string) {
		path := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("pkg/prod/prod.go", `package prod
import "github.com/ZoneCNH/testkitx/pkg/testkitx/golden"
var _ = golden.UpdateEnv
`)
	write("pkg/prod/prod_test.go", `package prod
import "github.com/ZoneCNH/testkitx/pkg/testkitx/assertx"
var _ = assertx.Failf
`)
	write("examples/example.go", `package examples
import "github.com/ZoneCNH/testkitx/pkg/testkitx/golden"
var _ = golden.UpdateEnv
`)

	violations, err := boundarytest.ScanProductionImports(root, "github.com/ZoneCNH/testkitx/pkg/testkitx/")
	if err != nil {
		t.Fatal(err)
	}
	if len(violations) != 1 || violations[0].ImportPath != "github.com/ZoneCNH/testkitx/pkg/testkitx/golden" {
		t.Fatalf("unexpected violations: %+v", violations)
	}
}
