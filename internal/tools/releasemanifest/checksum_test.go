package main

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestFileDigestRecordsPathAndSHA256(t *testing.T) {
	t.Parallel()

	path := t.TempDir() + "/contract.json"
	if err := os.WriteFile(path, []byte("abc"), 0o644); err != nil {
		t.Fatal(err)
	}

	digest, err := fileDigest(path)
	if err != nil {
		t.Fatal(err)
	}

	if digest.Path != path {
		t.Fatalf("path = %q, want %q", digest.Path, path)
	}
	const want = "sha256:ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
	if digest.SHA256 != want {
		t.Fatalf("sha256 = %q, want %q", digest.SHA256, want)
	}
}

func TestContractFilesCoverReleaseContracts(t *testing.T) {
	t.Parallel()

	required := []string{
		"contracts/config.schema.json",
		"contracts/docker-toolchain.schema.json",
		"contracts/downstream-adoption-proof.schema.json",
		"contracts/error.schema.json",
		"contracts/health.schema.json",
		"contracts/metrics.md",
	}
	seen := make(map[string]struct{}, len(contractFiles()))
	for _, path := range contractFiles() {
		seen[path] = struct{}{}
	}
	for _, path := range required {
		if _, ok := seen[path]; !ok {
			t.Fatalf("contractFiles() does not include %q", path)
		}
	}
}

func TestSourceDigestUsesTrackedFileNamesAndContents(t *testing.T) {
	repo := t.TempDir()
	runTestCommand(t, repo, "git", "init")

	files := map[string]string{
		"a.txt":          "alpha\n",
		"nested/b.txt":   "bravo\n",
		"nested/cfg.yml": "name: charlie\n",
	}
	for path, content := range files {
		fullPath := filepath.Join(repo, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	runTestCommand(t, repo, "git", "add", ".")
	chdir(t, repo)

	gotDigest, gotCount, err := sourceDigest()
	if err != nil {
		t.Fatal(err)
	}

	if gotCount != len(files) {
		t.Fatalf("tracked file count = %d, want %d", gotCount, len(files))
	}
	if want := expectedSourceDigest(files); gotDigest != want {
		t.Fatalf("source digest = %q, want %q", gotDigest, want)
	}
}

func TestModuleDigestsIncludesReplaceMetadata(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte(`module example.com/root

go 1.24

require example.com/dep v0.0.0

replace example.com/dep => ./dep
`), 0o644); err != nil {
		t.Fatal(err)
	}
	depDir := filepath.Join(root, "dep")
	if err := os.MkdirAll(depDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(depDir, "go.mod"), []byte("module example.com/dep\n\ngo 1.24\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("GOWORK", "off")
	chdir(t, root)

	modules, err := moduleDigests()
	if err != nil {
		t.Fatal(err)
	}

	var foundMain bool
	var foundReplace bool
	for _, module := range modules {
		if module.Path == "example.com/root" && module.Main {
			foundMain = true
		}
		if module.Path == "example.com/dep" && module.Replace != nil && module.Replace.Path == "./dep" {
			foundReplace = true
		}
	}
	if !foundMain {
		t.Fatalf("modules = %+v, want main module", modules)
	}
	if !foundReplace {
		t.Fatalf("modules = %+v, want replace metadata for example.com/dep", modules)
	}
}

func TestReadManifestChecksumRejectsEmptyFile(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "empty.json.sha256")
	os.WriteFile(path, []byte(""), 0o644)
	_, err := readManifestChecksum(path, "empty.json")
	if err == nil {
		t.Fatal("expected error for empty checksum file")
	}
}

func TestReadManifestChecksumRejectsShortDigest(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "short.json.sha256")
	os.WriteFile(path, []byte("abcd1234  short.json\n"), 0o644)
	_, err := readManifestChecksum(path, "short.json")
	if err == nil {
		t.Fatal("expected error for short digest")
	}
}

func TestReadManifestChecksumRejectsInvalidHex(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "bad.json.sha256")
	badHex := strings.Repeat("z0", 32)
	os.WriteFile(path, []byte(badHex+"  bad.json\n"), 0o644)
	_, err := readManifestChecksum(path, "bad.json")
	if err == nil {
		t.Fatal("expected error for invalid hex")
	}
}

func TestReadManifestChecksumRejectsWrongFilename(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "manifest.json.sha256")
	validHex := strings.Repeat("ab", 32)
	os.WriteFile(path, []byte(validHex+"  wrong.json\n"), 0o644)
	_, err := readManifestChecksum(path, "manifest.json")
	if err == nil {
		t.Fatal("expected error for wrong filename")
	}
}

func TestReadManifestChecksumAcceptsSHA256Prefix(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "prefixed.json.sha256")
	validHex := strings.Repeat("ab", 32)
	os.WriteFile(path, []byte("sha256:"+validHex+"  prefixed.json\n"), 0o644)
	got, err := readManifestChecksum(path, "prefixed.json")
	if err != nil {
		t.Fatal(err)
	}
	if got != validHex {
		t.Fatalf("checksum = %q, want %q", got, validHex)
	}
}

func TestReadManifestChecksumRejectsNonexistentFile(t *testing.T) {
	t.Parallel()
	_, err := readManifestChecksum("/nonexistent/path.json.sha256", "path.json")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestWriteManifestCreatesDirectory(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "sub", "dir", "manifest.json")
	manifest := Manifest{Module: "test", Version: "v0.1.0"}
	if err := writeManifest(path, manifest); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"module": "test"`) {
		t.Fatalf("manifest content = %q, want module field", string(data))
	}
}

func TestContainsReturnsFalseForMissingValue(t *testing.T) {
	t.Parallel()
	if contains([]string{"a", "b", "c"}, "d") {
		t.Fatal("contains([a,b,c], d) = true, want false")
	}
}

func TestContainsReturnsTrueForPresentValue(t *testing.T) {
	t.Parallel()
	if !contains([]string{"a", "b", "c"}, "b") {
		t.Fatal("contains([a,b,c], b) = false, want true")
	}
}

func TestRequireNonEmptyAppendsWhenEmpty(t *testing.T) {
	t.Parallel()
	var failures []string
	requireNonEmpty(&failures, "field", "  ")
	if len(failures) != 1 || !strings.Contains(failures[0], "field is required") {
		t.Fatalf("failures = %v, want 'field is required'", failures)
	}
}

func TestRequireNonEmptySkipsWhenNonEmpty(t *testing.T) {
	t.Parallel()
	var failures []string
	requireNonEmpty(&failures, "field", "value")
	if len(failures) != 0 {
		t.Fatalf("failures = %v, want empty", failures)
	}
}

func TestFirstLineHandlesMultiline(t *testing.T) {
	t.Parallel()
	got := firstLine("first\nsecond\nthird")
	if got != "first" {
		t.Fatalf("firstLine = %q, want %q", got, "first")
	}
}

func TestFirstLineHandlesSingleLine(t *testing.T) {
	t.Parallel()
	got := firstLine("only")
	if got != "only" {
		t.Fatalf("firstLine = %q, want %q", got, "only")
	}
}

func TestEnvDefaultReturnsEnvValue(t *testing.T) {
	t.Setenv("TESTKITX_TEST_ENV", "from_env")
	got := envDefault("TESTKITX_TEST_ENV", "fallback")
	if got != "from_env" {
		t.Fatalf("envDefault = %q, want %q", got, "from_env")
	}
}

func TestEnvDefaultReturnsFallback(t *testing.T) {
	t.Parallel()
	got := envDefault("TESTKITX_UNSET_VAR_xyz", "fallback")
	if got != "fallback" {
		t.Fatalf("envDefault = %q, want %q", got, "fallback")
	}
}

func expectedSourceDigest(files map[string]string) string {
	paths := make([]string, 0, len(files))
	for path := range files {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	digest := sha256.New()
	for _, path := range paths {
		sum := sha256.Sum256([]byte(files[path]))
		digest.Write([]byte(path))
		digest.Write([]byte{0})
		digest.Write([]byte(hex.EncodeToString(sum[:])))
		digest.Write([]byte{0})
	}
	return "sha256:" + hex.EncodeToString(digest.Sum(nil))
}
