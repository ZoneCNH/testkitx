package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
)

func verifyManifest(path string, requirePassed bool, requireClean bool, expectVersion string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var got Manifest
	if err := json.Unmarshal(data, &got); err != nil {
		return err
	}

	current, err := buildManifest()
	if err != nil {
		return err
	}

	var failures []string
	requireNonEmpty(&failures, "module", got.Module)
	requireNonEmpty(&failures, "version", got.Version)
	requireNonEmpty(&failures, "commit", got.Commit)
	requireNonEmpty(&failures, "tree_sha", got.TreeSHA)
	requireNonEmpty(&failures, "source_digest", got.SourceDigest)
	requireNonEmpty(&failures, "go_version", got.GoVersion)
	requireNonEmpty(&failures, "generated_at", got.GeneratedAt)
	requireNonEmpty(&failures, "generated_by", got.GeneratedBy)
	requireNonEmpty(&failures, "tree_state", got.TreeState)

	expectVersion = strings.TrimSpace(expectVersion)
	if _, err := time.Parse(time.RFC3339, got.GeneratedAt); err != nil {
		failures = append(failures, "generated_at must be RFC3339")
	}
	if expectVersion != "" && got.Version != expectVersion {
		failures = append(failures, fmt.Sprintf("version mismatch: got %q, want %q", got.Version, expectVersion))
	}
	if got.Module != current.Module {
		failures = append(failures, fmt.Sprintf("module mismatch: got %q, want %q", got.Module, current.Module))
	}
	if got.Commit != current.Commit {
		failures = append(failures, fmt.Sprintf("commit mismatch: got %q, want %q", got.Commit, current.Commit))
	}
	if got.TreeSHA != current.TreeSHA {
		failures = append(failures, fmt.Sprintf("tree_sha mismatch: got %q, want %q", got.TreeSHA, current.TreeSHA))
	}
	if got.SourceDigest != current.SourceDigest {
		failures = append(failures, "source_digest does not match current tracked file contents")
	}
	if got.TrackedFileCount != current.TrackedFileCount {
		failures = append(failures, fmt.Sprintf("tracked_file_count mismatch: got %d, want %d", got.TrackedFileCount, current.TrackedFileCount))
	}
	if got.TreeState != current.TreeState {
		failures = append(failures, fmt.Sprintf("tree_state mismatch: got %q, want %q", got.TreeState, current.TreeState))
	}
	if requireClean && got.TreeState != "clean" {
		failures = append(failures, fmt.Sprintf("tree_state must be clean, got %q", got.TreeState))
	}
	if !reflect.DeepEqual(got.Contracts, current.Contracts) {
		failures = append(failures, "contract fingerprints do not match current contract files")
	}
	if !reflect.DeepEqual(got.Dependencies, current.Dependencies) {
		failures = append(failures, "dependency inventory does not match go list -m -json all")
	}
	if !contains(got.Artifacts, "release/manifest/latest.json") {
		failures = append(failures, "artifacts must include release/manifest/latest.json")
	}
	if !contains(got.Artifacts, "release/manifest/latest.json.sha256") {
		failures = append(failures, "artifacts must include release/manifest/latest.json.sha256")
	}
	if err := verifyManifestChecksum(path, data); err != nil {
		failures = append(failures, err.Error())
	}
	if got.Tools["go"] == "" {
		failures = append(failures, "tools.go must be recorded")
	}
	failures = append(failures, validateChecks(got.Checks, requirePassed)...)

	if len(failures) > 0 {
		return errors.New("release evidence verification failed:\n - " + strings.Join(failures, "\n - "))
	}
	return nil
}

func validateChecks(checks map[string]string, requirePassed bool) []string {
	var failures []string
	for _, name := range checkNames() {
		status := strings.TrimSpace(checks[name])
		if status == "" {
			failures = append(failures, "checks."+name+" is required")
			continue
		}
		if requirePassed && status != "passed" {
			failures = append(failures, fmt.Sprintf("checks.%s must be passed, got %q", name, status))
		}
	}
	return failures
}
