package contracts

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx"
)

func TestReleaseReadinessDocsAndMetadataStayAligned(t *testing.T) {
	t.Parallel()

	version := testkitx.Version
	if version == "" {
		t.Fatal("testkitx.Version must not be empty")
	}

	repoContract := readFile(t, "../.repo-contract.yaml")
	changelog := readFile(t, "../CHANGELOG.md")
	features := readFile(t, "../FEATURES.md")
	acceptance := readFile(t, "../ACCEPTANCE.md")

	latestTag := mustFindVersion(t, repoContract, `(?m)^\s*latest_git_tag:\s*(v\d+\.\d+\.\d+)\s*$`)
	if latestTag != version {
		t.Fatalf(".repo-contract.yaml latest_git_tag = %q, want %q", latestTag, version)
	}

	wantHeading := "## " + version + " - "
	if !strings.Contains(changelog, wantHeading) {
		t.Fatalf("CHANGELOG.md missing %q", wantHeading)
	}

	for _, doc := range []struct {
		name string
		text string
	}{
		{name: "FEATURES.md", text: features},
		{name: "ACCEPTANCE.md", text: acceptance},
	} {
		if !strings.Contains(doc.text, version) {
			t.Fatalf("%s does not reference current version %q", doc.name, version)
		}
	}

	if !strings.Contains(features, "L1 测试专用能力库") {
		t.Fatalf("FEATURES.md should describe the current library identity")
	}
	if !strings.Contains(acceptance, "release-final-check") {
		t.Fatalf("ACCEPTANCE.md should describe the release-final-check gate")
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

func mustFindVersion(t *testing.T, text, pattern string) string {
	t.Helper()

	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(text)
	if len(match) != 2 {
		t.Fatalf("pattern %q not found", pattern)
	}
	return match[1]
}
