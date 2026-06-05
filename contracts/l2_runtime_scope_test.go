package contracts

import (
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestL2RuntimeScopeExcludesProvidersAndReleaseDecisions(t *testing.T) {
	t.Parallel()

	files := l2RuntimeSourceFiles(t)
	modulePath := l2RuntimeModulePath(t)
	if len(files) == 0 {
		t.Log("no L2 runtime source files found yet; scope guard will apply when packages are added")
	}

	for _, path := range files {
		path := path
		t.Run(path, func(t *testing.T) {
			t.Parallel()

			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read %s: %v", path, err)
			}
			text := string(content)
			for _, rule := range l2RuntimeForbiddenScopeRules() {
				if rule.pattern.MatchString(text) {
					t.Fatalf("%s violates L2 runtime scope rule %q", path, rule.name)
				}
			}
			for _, importPath := range l2RuntimeImports(t, path, content) {
				if l2RuntimeForbiddenRepositoryImport(importPath, modulePath) {
					t.Fatalf("%s violates L2 runtime scope rule %q via import %q", path, "L2 repository dependency", importPath)
				}
			}
		})
	}
}

type forbiddenScopeRule struct {
	name    string
	pattern *regexp.Regexp
}

func l2RuntimeForbiddenScopeRules() []forbiddenScopeRule {
	return []forbiddenScopeRule{
		{
			name:    "provider-specific dependency or implementation detail",
			pattern: regexp.MustCompile(`(?i)(github\.com/(redis|segmentio|shopify|confluentinc|jackc|lib/pq|minio|aliyun|IBM/sarama|ClickHouse|taosdata)|\b(redis|kafka|postgres|postgresql|pgx|sarama|minio|aliyun|oss|clickhouse|taos|dynamodb|pubsubx|redisx|kafkax|sqlx)\b)`),
		},
		{
			name:    "external service connection workflow",
			pattern: regexp.MustCompile(`(?i)\b(dial(?:context)?|connect(?:context)?|open\s*\(|new\s*(redis|kafka|postgres|sql|pubsub|object|timeseries|columnstore)\s*client)\b`),
		},
		{
			name:    "release decision or release-readiness logic",
			pattern: regexp.MustCompile(`(?i)\b(release[-_ ]?(decision|readiness|ready|gate)|can[-_ ]?release|approve[-_ ]?release|block[-_ ]?release|go[-_ ]?no[-_ ]?go)\b`),
		},
		{
			name:    "layer or secret scan workflow",
			pattern: regexp.MustCompile(`(?i)\b(layer[-_ ]?scan|scan[-_ ]?layer|secret[-_ ]?scan|scan[-_ ]?secret)\b`),
		},
	}
}

func l2RuntimeModulePath(t *testing.T) string {
	t.Helper()

	data, err := os.ReadFile("../go.mod")
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[0] == "module" {
			return fields[1]
		}
	}
	t.Fatal("go.mod module path not found")
	return ""
}

func l2RuntimeImports(t *testing.T, path string, content []byte) []string {
	t.Helper()

	file, err := parser.ParseFile(token.NewFileSet(), path, content, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("parse imports from %s: %v", path, err)
	}
	imports := make([]string, 0, len(file.Imports))
	for _, importSpec := range file.Imports {
		importPath, err := strconv.Unquote(importSpec.Path.Value)
		if err != nil {
			t.Fatalf("parse import path in %s: %v", path, err)
		}
		imports = append(imports, importPath)
	}
	return imports
}

func l2RuntimeForbiddenRepositoryImport(importPath, modulePath string) bool {
	if importPath == modulePath || strings.HasPrefix(importPath, modulePath+"/") {
		return false
	}

	lowerImport := strings.ToLower(importPath)
	for _, forbidden := range []string{
		"github.com/zonecnh/redisx",
		"github.com/zonecnh/foundationx",
		"github.com/zonecnh/kafkax",
		"github.com/zonecnh/pubsubx",
		"github.com/zonecnh/servicex",
		"github.com/zonecnh/xlib",
		"github.com/bytechainx/redisx",
		"github.com/bytechainx/foundationx",
		"github.com/bytechainx/kafkax",
		"github.com/bytechainx/pubsubx",
		"github.com/bytechainx/servicex",
		"github.com/bytechainx/xlib",
	} {
		if lowerImport == forbidden || strings.HasPrefix(lowerImport, forbidden+"/") {
			return true
		}
	}
	return false
}

func l2RuntimeSourceFiles(t *testing.T) []string {
	t.Helper()

	roots := []string{
		"../pkg/testkitx/requirex",
		"../pkg/testkitx/contract",
		"../pkg/testkitx/evidence",
		"../pkg/testkitx/servicex",
		"../requirex",
		"../contract",
		"../evidence",
		"../servicex",
	}
	files := make([]string, 0)
	for _, root := range roots {
		info, err := os.Stat(root)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			t.Fatalf("stat %s: %v", root, err)
		}
		if !info.IsDir() {
			t.Fatalf("L2 runtime root %s is not a directory", root)
		}
		if err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				name := d.Name()
				if strings.HasPrefix(name, ".") || name == "testdata" {
					return filepath.SkipDir
				}
				return nil
			}
			if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
				files = append(files, path)
			}
			return nil
		}); err != nil {
			t.Fatalf("walk %s: %v", root, err)
		}
	}
	return files
}
