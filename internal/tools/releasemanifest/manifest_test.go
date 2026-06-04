package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestBuildChecksUsesGlobalAndSpecificStatus(t *testing.T) {
	t.Setenv("CHECK_STATUS", "passed")
	t.Setenv("LINT_STATUS", "failed")

	checks := buildChecks()

	if checks["fmt"] != "passed" {
		t.Fatalf("fmt status = %q, want passed", checks["fmt"])
	}
	if checks["lint"] != "failed" {
		t.Fatalf("lint status = %q, want failed", checks["lint"])
	}
}

func TestBuildManifestRecordsCurrentRepositoryFacts(t *testing.T) {
	t.Setenv("GOWORK", "off")
	t.Setenv("VERSION", "v9.9.9-test")
	t.Setenv("GENERATED_BY", "releasemanifest-test")
	t.Setenv("CHECK_STATUS", "passed")
	chdir(t, repoRoot(t))

	manifest, err := buildManifest()
	if err != nil {
		t.Fatal(err)
	}

	if manifest.Module != "github.com/ZoneCNH/testkitx" {
		t.Fatalf("module = %q, want github.com/ZoneCNH/testkitx", manifest.Module)
	}
	if manifest.Version != "v9.9.9-test" {
		t.Fatalf("version = %q, want v9.9.9-test", manifest.Version)
	}
	if manifest.GeneratedBy != "releasemanifest-test" {
		t.Fatalf("generated_by = %q, want releasemanifest-test", manifest.GeneratedBy)
	}
	if _, err := time.Parse(time.RFC3339, manifest.GeneratedAt); err != nil {
		t.Fatalf("generated_at = %q, want RFC3339: %v", manifest.GeneratedAt, err)
	}
	if !strings.HasPrefix(manifest.SourceDigest, "sha256:") {
		t.Fatalf("source_digest = %q, want sha256 prefix", manifest.SourceDigest)
	}
	if manifest.TrackedFileCount == 0 {
		t.Fatal("tracked_file_count = 0, want tracked files")
	}
	if len(manifest.Contracts) != len(contractFiles()) {
		t.Fatalf("len(contracts) = %d, want %d", len(manifest.Contracts), len(contractFiles()))
	}
	if len(manifest.Dependencies) == 0 || manifest.Dependencies[0].Path != manifest.Module || !manifest.Dependencies[0].Main {
		t.Fatalf("dependencies[0] = %+v, want main module %q", manifest.Dependencies, manifest.Module)
	}
	if manifest.Tools["go"] == "" {
		t.Fatal("tools.go is empty")
	}
	if !contains(manifest.Artifacts, "release/manifest/latest.json") {
		t.Fatalf("artifacts = %v, want release/manifest/latest.json", manifest.Artifacts)
	}
	if !contains(manifest.Artifacts, "release/manifest/latest.json.sha256") {
		t.Fatalf("artifacts = %v, want release/manifest/latest.json.sha256", manifest.Artifacts)
	}
	for _, name := range checkNames() {
		if manifest.Checks[name] != "passed" {
			t.Fatalf("checks[%q] = %q, want passed", name, manifest.Checks[name])
		}
	}
	if manifest.TreeState != "clean" && manifest.TreeState != "dirty" {
		t.Fatalf("tree_state = %q, want clean or dirty", manifest.TreeState)
	}
}

func TestWriteManifestCreatesParentAndWritesIndentedJSON(t *testing.T) {
	manifest := Manifest{
		Module:           "example.com/lib",
		Version:          "v1.2.3",
		Commit:           "abc123",
		TreeSHA:          "tree123",
		SourceDigest:     "sha256:source",
		TrackedFileCount: 1,
		GoVersion:        "go1.23.0",
		GeneratedAt:      "2026-01-02T03:04:05Z",
		GeneratedBy:      "test",
		TreeState:        "clean",
		Checks:           map[string]string{"fmt": "passed"},
		Tools:            map[string]string{"go": "go version"},
		Artifacts:        []string{"release/manifest/latest.json", "release/manifest/latest.json.sha256"},
	}
	path := filepath.Join(t.TempDir(), "release", "manifest", "latest.json")

	if err := writeManifest(path, manifest); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(data) {
		t.Fatalf("manifest JSON is invalid: %s", data)
	}
	if !strings.Contains(string(data), "\n  ") {
		t.Fatalf("manifest JSON is not indented: %s", data)
	}
	assertManifestChecksum(t, path)

	var got Manifest
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Module != manifest.Module || got.Version != manifest.Version {
		t.Fatalf("round-trip manifest = %+v, want %+v", got, manifest)
	}
}
