package main

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestFileDigestRecordsPathAndSHA256(t *testing.T) {
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

go 1.23

require example.com/dep v0.0.0

replace example.com/dep => ./dep
`), 0o644); err != nil {
		t.Fatal(err)
	}
	depDir := filepath.Join(root, "dep")
	if err := os.MkdirAll(depDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(depDir, "go.mod"), []byte("module example.com/dep\n\ngo 1.23\n"), 0o644); err != nil {
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
