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
	if violations[0].Line == 0 {
		t.Fatalf("expected Line to be set, got 0")
	}
	if violations[0].Reason == "" {
		t.Fatalf("expected Reason to be set")
	}
}

func TestScanWithAllowedImports(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	path := filepath.Join(root, "pkg", "app", "app.go")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`package app
import "github.com/ZoneCNH/testkitx/pkg/testkitx/golden"
var _ = golden.UpdateEnv
`), 0o644); err != nil {
		t.Fatal(err)
	}

	// Without whitelist — should find violation.
	v, err := boundarytest.Scan(boundarytest.ScanConfig{
		Dir:             root,
		ForbiddenPrefix: "github.com/ZoneCNH/testkitx/pkg/testkitx/",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}

	// With whitelist — violation should be suppressed.
	v, err = boundarytest.Scan(boundarytest.ScanConfig{
		Dir:             root,
		ForbiddenPrefix: "github.com/ZoneCNH/testkitx/pkg/testkitx/",
		AllowedImports:  []string{"github.com/ZoneCNH/testkitx/pkg/testkitx/golden"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 0 {
		t.Fatalf("expected 0 violations with whitelist, got %d: %+v", len(v), v)
	}
}

func TestViolationStringFormatting(t *testing.T) {
	t.Parallel()
	v := boundarytest.Violation{
		File:       "/src/app.go",
		Line:       5,
		ImportPath: "github.com/ZoneCNH/testkitx/pkg/testkitx/golden",
		Reason:     "production code imports testkitx internal package",
	}
	s := v.String()
	if s == "" {
		t.Fatal("expected non-empty string")
	}
}
