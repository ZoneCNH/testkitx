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

func TestScanNonexistentDir(t *testing.T) {
	t.Parallel()
	v, err := boundarytest.Scan(boundarytest.ScanConfig{
		Dir:             "/nonexistent/dir/path",
		ForbiddenPrefix: "testkitx",
	})
	if err == nil {
		if len(v) != 0 {
			t.Fatalf("expected no violations for nonexistent dir, got %d", len(v))
		}
	}
}

func TestScanEmptyDir(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	v, err := boundarytest.Scan(boundarytest.ScanConfig{
		Dir:             dir,
		ForbiddenPrefix: "testkitx",
	})
	if err != nil {
		t.Fatalf("Scan empty dir: %v", err)
	}
	if len(v) != 0 {
		t.Fatalf("expected no violations in empty dir, got %d", len(v))
	}
}

func TestScanTestFilesExcluded(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	path := filepath.Join(root, "pkg", "prod", "prod_test.go")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`package prod
import "github.com/ZoneCNH/testkitx/pkg/testkitx/golden"
var _ = golden.UpdateEnv
`), 0o644); err != nil {
		t.Fatal(err)
	}
	v, err := boundarytest.Scan(boundarytest.ScanConfig{
		Dir:             root,
		ForbiddenPrefix: "github.com/ZoneCNH/testkitx/pkg/testkitx/",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 0 {
		t.Fatalf("test files should be excluded, got %d violations: %+v", len(v), v)
	}
}

func TestScanTestkitDirectoryExcluded(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	path := filepath.Join(root, "testkit", "helper.go")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`package testkit
import "github.com/ZoneCNH/testkitx/pkg/testkitx/golden"
var _ = golden.UpdateEnv
`), 0o644); err != nil {
		t.Fatal(err)
	}
	v, err := boundarytest.Scan(boundarytest.ScanConfig{
		Dir:             root,
		ForbiddenPrefix: "github.com/ZoneCNH/testkitx/pkg/testkitx/",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 0 {
		t.Fatalf("testkit/ directory should be excluded, got %d violations", len(v))
	}
}

func TestScanNoMatchingImports(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	path := filepath.Join(root, "pkg", "app", "app.go")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`package app
import "fmt"
var _ = fmt.Println
`), 0o644); err != nil {
		t.Fatal(err)
	}
	v, err := boundarytest.Scan(boundarytest.ScanConfig{
		Dir:             root,
		ForbiddenPrefix: "github.com/ZoneCNH/testkitx/",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 0 {
		t.Fatalf("no forbidden imports, got %d violations", len(v))
	}
}

func TestScanMultipleViolations(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	for _, f := range []struct {
		path, body string
	}{
		{"pkg/a/a.go", "package a\nimport \"github.com/ZoneCNH/testkitx/pkg/testkitx/golden\"\nvar _ = golden.UpdateEnv\n"},
		{"pkg/b/b.go", "package b\nimport \"github.com/ZoneCNH/testkitx/pkg/testkitx/assertx\"\nvar _ = assertx.Failf\n"},
	} {
		p := filepath.Join(root, f.path)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(f.body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	v, err := boundarytest.Scan(boundarytest.ScanConfig{
		Dir:             root,
		ForbiddenPrefix: "github.com/ZoneCNH/testkitx/pkg/testkitx/",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 2 {
		t.Fatalf("expected 2 violations, got %d: %+v", len(v), v)
	}
}

func TestScanLegacyWrapperEmptyDir(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	v, err := boundarytest.ScanProductionImports(dir, "testkitx")
	if err != nil {
		t.Fatalf("ScanProductionImports: %v", err)
	}
	if len(v) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(v))
	}
}
