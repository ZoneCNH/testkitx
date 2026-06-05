package contract_test

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCapabilityRunnersAvoidProviderSpecificImports(t *testing.T) {
	t.Parallel()
	families := []string{"kv", "sql", "pubsub", "eventlog", "objectstore", "columnstore", "timeseries"}
	blocked := []string{"redis", "kafka", "postgres", "pq", "database/sql"}

	for _, family := range families {
		family := family
		t.Run(family, func(t *testing.T) {
			t.Parallel()
			files, err := os.ReadDir(family)
			if err != nil {
				t.Fatalf("read %s: %v", family, err)
			}
			for _, file := range files {
				if file.IsDir() || !strings.HasSuffix(file.Name(), ".go") || strings.HasSuffix(file.Name(), "_test.go") {
					continue
				}
				path := filepath.Join(family, file.Name())
				parsed, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
				if err != nil {
					t.Fatalf("parse imports from %s: %v", path, err)
				}
				for _, spec := range parsed.Imports {
					importPath := strings.Trim(spec.Path.Value, "\"")
					for _, fragment := range blocked {
						if strings.Contains(importPath, fragment) {
							t.Fatalf("%s imports provider-specific package %q", path, importPath)
						}
					}
				}
			}
		})
	}
}
