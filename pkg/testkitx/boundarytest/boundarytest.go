// Package boundarytest scans for illegal production imports of testkitx.
package boundarytest

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
)

type Violation struct {
	File       string
	ImportPath string
}

func ScanProductionImports(root, forbiddenPrefix string) ([]Violation, error) {
	forbiddenPrefix = strings.TrimSuffix(forbiddenPrefix, "/")
	var out []Violation
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		if strings.HasPrefix(rel, "testkit"+string(filepath.Separator)) || strings.HasPrefix(rel, "tools"+string(filepath.Separator)) || strings.HasPrefix(rel, "examples"+string(filepath.Separator)) {
			return nil
		}
		file, parseErr := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if parseErr != nil {
			return parseErr
		}
		for _, imp := range file.Imports {
			p := strings.Trim(imp.Path.Value, "\"")
			if p == forbiddenPrefix || strings.HasPrefix(p, forbiddenPrefix+"/") {
				out = append(out, Violation{File: path, ImportPath: p})
			}
		}
		return nil
	})
	return out, err
}
