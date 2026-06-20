// Package boundarytest scans for illegal production imports of testkitx.
package boundarytest

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
)

// Violation describes a single production-import boundary violation.
type Violation struct {
	File       string
	Line       int
	ImportPath string
	Reason     string
}

func (v Violation) String() string {
	return fmt.Sprintf("%s:%d: %s — %s", v.File, v.Line, v.ImportPath, v.Reason)
}

// ScanConfig controls what Scan considers a violation.
type ScanConfig struct {
	// Dir is the root directory to walk.
	Dir string
	// ForbiddenPrefix is the import prefix that production code must not use.
	ForbiddenPrefix string
	// AllowedImports is a set of exact import paths that are always permitted
	// even when they match ForbiddenPrefix (e.g. shared error types).
	AllowedImports []string
}

// Scan walks Dir and returns violations for production .go files that import
// ForbiddenPrefix. _test.go files, and directories under testkit/, tools/,
// examples/ are excluded. AllowedImports is consulted as a whitelist.
func Scan(cfg ScanConfig) ([]Violation, error) {
	prefix := strings.TrimSuffix(cfg.ForbiddenPrefix, "/")
	allowed := make(map[string]bool, len(cfg.AllowedImports))
	for _, a := range cfg.AllowedImports {
		allowed[a] = true
	}
	var out []Violation
	fset := token.NewFileSet()
	err := filepath.WalkDir(cfg.Dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return err
		}
		rel, _ := filepath.Rel(cfg.Dir, path)
		for _, skip := range []string{"testkit", "tools", "examples"} {
			if strings.HasPrefix(rel, skip+string(filepath.Separator)) {
				return nil
			}
		}
		file, parseErr := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if parseErr != nil {
			return parseErr
		}
		for _, imp := range file.Imports {
			p := strings.Trim(imp.Path.Value, "\"")
			if p != prefix && !strings.HasPrefix(p, prefix+"/") {
				continue
			}
			if allowed[p] {
				continue
			}
			fileInfo := fset.File(imp.Pos())
			line := 0
			if fileInfo != nil {
				line = fileInfo.Line(imp.Pos())
			}
			out = append(out, Violation{
				File:       path,
				Line:       line,
				ImportPath: p,
				Reason:     "production code imports testkitx internal package",
			})
		}
		return nil
	})
	return out, err
}

// ScanProductionImports is the legacy convenience wrapper kept for backward
// compatibility. Prefer Scan with ScanConfig for new code.
func ScanProductionImports(root, forbiddenPrefix string) ([]Violation, error) {
	return Scan(ScanConfig{Dir: root, ForbiddenPrefix: forbiddenPrefix})
}

// T is the subset of testing.TB used by boundary assertions.
type T interface {
	Helper()
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
}

// BoundaryCheck scans the given module directory for production imports of
// testkitx and fails t if any are found. Per SPEC FR-009.
//
// The module argument should be the root directory of the module to scan.
// Pass "." to scan the current module.
func BoundaryCheck(tt T, module string) {
	tt.Helper()
	violations, err := Scan(ScanConfig{
		Dir:             module,
		ForbiddenPrefix: "github.com/ZoneCNH/testkitx",
	})
	if err != nil {
		tt.Fatalf("boundary check failed: %v", err)
	}
	for _, v := range violations {
		tt.Errorf("boundary violation: %s", v.String())
	}
	if len(violations) > 0 {
		tt.Errorf("BoundaryCheck: %d production import(s) of testkitx detected", len(violations))
	}
}
