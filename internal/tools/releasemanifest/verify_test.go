package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateChecksRequiresPassedStatuses(t *testing.T) {
	t.Parallel()

	checks := make(map[string]string, len(checkNames()))
	for _, name := range checkNames() {
		checks[name] = "passed"
	}
	checks["security"] = "unknown"

	failures := validateChecks(checks, true)

	if len(failures) != 1 {
		t.Fatalf("len(failures) = %d, want 1: %v", len(failures), failures)
	}
	if !strings.Contains(failures[0], "checks.security") {
		t.Fatalf("failure = %q, want security check failure", failures[0])
	}
}

func TestValidateChecksRejectsInvalidStatusWithoutRequirePassed(t *testing.T) {
	t.Parallel()

	checks := make(map[string]string, len(checkNames()))
	for _, name := range checkNames() {
		checks[name] = "unknown"
	}
	checks["security"] = "bogus"

	failures := validateChecks(checks, false)

	if len(failures) != 1 {
		t.Fatalf("len(failures) = %d, want 1: %v", len(failures), failures)
	}
	if !strings.Contains(failures[0], `checks.security has invalid status "bogus"`) {
		t.Fatalf("failure = %q, want invalid status failure", failures[0])
	}
}

func TestVerifyManifestAcceptsFreshManifestAndRejectsDrift(t *testing.T) {
	t.Setenv("GOWORK", "off")
	t.Setenv("CHECK_STATUS", "passed")
	chdir(t, repoRoot(t))

	manifest, err := buildManifest()
	if err != nil {
		t.Fatal(err)
	}

	goodPath := filepath.Join(t.TempDir(), "latest.json")
	if err := writeManifest(goodPath, manifest); err != nil {
		t.Fatal(err)
	}
	if err := verifyManifest(goodPath, true, false, ""); err != nil {
		t.Fatalf("verify fresh manifest: %v", err)
	}

	manifest.SourceDigest = "sha256:bad"
	manifest.Checks["lint"] = "unknown"
	badPath := filepath.Join(t.TempDir(), "stale.json")
	if err := writeManifest(badPath, manifest); err != nil {
		t.Fatal(err)
	}

	err = verifyManifest(badPath, true, false, "")
	if err == nil {
		t.Fatal("verify stale manifest succeeded, want error")
	}
	message := err.Error()
	for _, want := range []string{
		"source_digest does not match current tracked file contents",
		`checks.lint must be passed, got "unknown"`,
	} {
		if !strings.Contains(message, want) {
			t.Fatalf("error = %q, want substring %q", message, want)
		}
	}
}

func TestValidateChecksRejectsEmptyStatus(t *testing.T) {
	t.Parallel()

	checks := make(map[string]string, len(checkNames()))
	for _, name := range checkNames() {
		checks[name] = "passed"
	}
	checks["security"] = ""

	failures := validateChecks(checks, false)

	if len(failures) != 1 {
		t.Fatalf("len(failures) = %d, want 1: %v", len(failures), failures)
	}
	if !strings.Contains(failures[0], "checks.security is required") {
		t.Fatalf("failure = %q, want 'checks.security is required'", failures[0])
	}
}

func TestValidateChecksAcceptsValidStatuses(t *testing.T) {
	t.Parallel()

	checks := make(map[string]string, len(checkNames()))
	for _, name := range checkNames() {
		checks[name] = "passed"
	}

	failures := validateChecks(checks, true)

	if len(failures) != 0 {
		t.Fatalf("failures = %v, want empty", failures)
	}
}

func TestValidCheckStatus(t *testing.T) {
	t.Parallel()
	for _, status := range []string{"passed", "failed", "skipped", "unknown"} {
		if !validCheckStatus(status) {
			t.Fatalf("validCheckStatus(%q) = false, want true", status)
		}
	}
	if validCheckStatus("bogus") {
		t.Fatal("validCheckStatus(bogus) = true, want false")
	}
}

func TestVerifyManifestRejectsNonexistentFile(t *testing.T) {
	t.Parallel()
	err := verifyManifest("/nonexistent/manifest.json", false, false, "")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestVerifyManifestRejectsInvalidJSON(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "bad.json")
	os.WriteFile(path, []byte("not json"), 0o644)
	err := verifyManifest(path, false, false, "")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestVerifyManifestRequiresCleanTree(t *testing.T) {
	t.Setenv("GOWORK", "off")
	t.Setenv("CHECK_STATUS", "passed")
	chdir(t, repoRoot(t))

	manifest, err := buildManifest()
	if err != nil {
		t.Fatal(err)
	}
	manifest.TreeState = "dirty"

	path := filepath.Join(t.TempDir(), "dirty.json")
	if err := writeManifest(path, manifest); err != nil {
		t.Fatal(err)
	}

	err = verifyManifest(path, true, true, "")
	if err == nil {
		t.Fatal("verify dirty manifest with requireClean succeeded, want error")
	}
	if !strings.Contains(err.Error(), `tree_state must be clean, got "dirty"`) {
		t.Fatalf("error = %q, want require-clean failure", err)
	}
}
