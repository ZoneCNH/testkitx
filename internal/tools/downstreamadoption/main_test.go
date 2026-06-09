package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ZoneCNH/testkitx/internal/cliutil"
)

func TestRunCLIVerifiesValidEvidence(t *testing.T) {
	path := writeEvidence(t, validEvidence())
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI("downstreamadoption", []string{"-verify", path}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("runCLI verify exit code = %d, want 0; stderr: %s", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	if want := "downstream adoption evidence verified: " + path; !strings.Contains(stdout.String(), want) {
		t.Fatalf("stdout = %q, want substring %q", stdout.String(), want)
	}
}

func TestRunCLIVerifyRejectsMissingArtifactURL(t *testing.T) {
	t.Parallel()
	evidence := validEvidence()
	evidence.ArtifactURL = ""
	path := writeEvidence(t, evidence)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI("downstreamadoption", []string{"-verify", path}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("runCLI verify exit code = %d, want 1; stdout: %s; stderr: %s", code, stdout.String(), stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if want := "artifact_url is required"; !strings.Contains(stderr.String(), want) {
		t.Fatalf("stderr = %q, want substring %q", stderr.String(), want)
	}
}

func TestRunCLIVerifyRejectsProductionImports(t *testing.T) {
	t.Parallel()
	evidence := validEvidence()
	evidence.ProductionImports = 1
	path := writeEvidence(t, evidence)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI("downstreamadoption", []string{"-verify", path}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("runCLI verify exit code = %d, want 1; stdout: %s; stderr: %s", code, stdout.String(), stderr.String())
	}
	if want := "production_imports must be 0, got 1"; !strings.Contains(stderr.String(), want) {
		t.Fatalf("stderr = %q, want substring %q", stderr.String(), want)
	}
}

func TestValidateEvidenceAcceptsNumericWorkflowRunID(t *testing.T) {
	t.Parallel()
	data := []byte(`{
  "status": "passed",
  "repository": "ZoneCNH/xlib-standard",
  "commit": "1111111111111111111111111111111111111111",
  "tree_sha": "2222222222222222222222222222222222222222",
  "workflow_run_id": 26929076767,
  "artifact_url": "https://github.com/ZoneCNH/xlib-standard/actions/runs/26929076767",
  "sha256": "sha256:3333333333333333333333333333333333333333333333333333333333333333",
  "test_imports": 1,
  "production_imports": 0,
  "commands": [{"command": "GOWORK=off make ci", "status": "passed"}],
  "gates": {"ci": "passed"}
}`)
	var evidence Evidence
	if err := json.Unmarshal(data, &evidence); err != nil {
		t.Fatal(err)
	}

	if failures := validateEvidence(evidence); len(failures) != 0 {
		t.Fatalf("validateEvidence failures = %v, want none", failures)
	}
}

func TestValidateEvidenceRejectsMissingRequiredFields(t *testing.T) {
	t.Parallel()
	failures := validateEvidence(Evidence{})

	for _, want := range []string{
		`status must be "passed"`,
		"repository is required",
		"commit must be 40 hex characters",
		"tree_sha must be 40 hex characters",
		"workflow_run_id is required",
		"artifact_url is required",
		"sha256 must be 64 hex characters",
		"test_imports must be greater than 0",
		`commands must include passed command "GOWORK=off make ci"`,
		"gates must include at least one gate status",
	} {
		if !containsFailure(failures, want) {
			t.Fatalf("failures = %v, want substring %q", failures, want)
		}
	}
}

func TestValidateEvidenceRejectsPrefixedGitSHA(t *testing.T) {
	t.Parallel()
	evidence := validEvidence()
	evidence.Commit = "sha256:1111111111111111111111111111111111111111"

	failures := validateEvidence(evidence)

	if !containsFailure(failures, "commit must be 40 hex characters") {
		t.Fatalf("failures = %v, want prefixed commit rejection", failures)
	}
}

func TestValidateEvidenceRejectsMissingGates(t *testing.T) {
	t.Parallel()
	evidence := validEvidence()
	evidence.Gates = nil

	failures := validateEvidence(evidence)

	if !containsFailure(failures, "gates must include at least one gate status") {
		t.Fatalf("failures = %v, want missing gates rejection", failures)
	}
}

func TestRunCLIRequiresVerifyFlag(t *testing.T) {
	t.Parallel()
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI("downstreamadoption", nil, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("runCLI exit code = %d, want 1; stdout: %s; stderr: %s", code, stdout.String(), stderr.String())
	}
	if want := "-verify is required"; !strings.Contains(stderr.String(), want) {
		t.Fatalf("stderr = %q, want substring %q", stderr.String(), want)
	}
}

func TestRunCLIRejectsUnknownFlag(t *testing.T) {
	t.Parallel()
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI("downstreamadoption", []string{"-unknown"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("runCLI unknown flag exit code = %d, want 2; stderr: %s", code, stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if want := "flag provided but not defined"; !strings.Contains(stderr.String(), want) {
		t.Fatalf("stderr = %q, want substring %q", stderr.String(), want)
	}
}

func TestPrintCLIMessageReportsWriterFailure(t *testing.T) {
	t.Parallel()
	if code := cliutil.PrintCLIStatus(errorWriter{}, "ok\n"); code != 1 {
		t.Fatalf("PrintCLIStatus exit code = %d, want 1", code)
	}
	if code := cliutil.PrintCLIError(errorWriter{}, errors.New("boom")); code != 1 {
		t.Fatalf("PrintCLIError exit code = %d, want 1", code)
	}
}

func validEvidence() Evidence {
	return Evidence{
		Status:            "passed",
		Repository:        "ZoneCNH/xlib-standard",
		Commit:            "1111111111111111111111111111111111111111",
		TreeSHA:           "2222222222222222222222222222222222222222",
		WorkflowRunID:     flexibleString("26929076767"),
		WorkflowURL:       "https://github.com/ZoneCNH/xlib-standard/actions/runs/26929076767",
		ArtifactURL:       "https://github.com/ZoneCNH/xlib-standard/actions/runs/26929076767",
		SHA256:            "sha256:3333333333333333333333333333333333333333333333333333333333333333",
		TestImports:       1,
		ProductionImports: 0,
		Commands: []CommandEvidence{
			{Command: "GOWORK=off make ci", Status: "passed"},
		},
		Gates: map[string]string{
			"ci": "passed",
		},
	}
}

func writeEvidence(t *testing.T, evidence Evidence) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "downstream-adoption.json")
	data, err := json.MarshalIndent(evidence, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func containsFailure(failures []string, want string) bool {
	for _, failure := range failures {
		if strings.Contains(failure, want) {
			return true
		}
	}
	return false
}

func TestRequireNumericRejectsEmptyValue(t *testing.T) {
	t.Parallel()
	var failures []string
	requireNumeric(&failures, "workflow_run_id", "  ")
	if len(failures) != 1 || !strings.Contains(failures[0], "is required") {
		t.Fatalf("failures = %v, want 'is required'", failures)
	}
}

func TestRequireNumericRejectsNonDigits(t *testing.T) {
	t.Parallel()
	var failures []string
	requireNumeric(&failures, "workflow_run_id", "abc123")
	if len(failures) != 1 || !strings.Contains(failures[0], "digits only") {
		t.Fatalf("failures = %v, want 'digits only'", failures)
	}
}

func TestRequireHTTPURLRejectsEmptyValue(t *testing.T) {
	t.Parallel()
	var failures []string
	requireHTTPURL(&failures, "artifact_url", "  ")
	if len(failures) != 1 || !strings.Contains(failures[0], "is required") {
		t.Fatalf("failures = %v, want 'is required'", failures)
	}
}

func TestRequireHTTPURLRejectsInvalidURL(t *testing.T) {
	t.Parallel()
	var failures []string
	requireHTTPURL(&failures, "artifact_url", "://bad")
	if len(failures) != 1 || !strings.Contains(failures[0], "absolute URL") {
		t.Fatalf("failures = %v, want 'absolute URL'", failures)
	}
}

func TestRequireHTTPURLRejectsNonHTTPScheme(t *testing.T) {
	t.Parallel()
	var failures []string
	requireHTTPURL(&failures, "artifact_url", "ftp://example.com/file")
	if len(failures) != 1 || !strings.Contains(failures[0], "http or https") {
		t.Fatalf("failures = %v, want 'http or https'", failures)
	}
}

func TestRequireHexStringRejectsWrongLength(t *testing.T) {
	t.Parallel()
	var failures []string
	requireHexString(&failures, "commit", "abcd", 20)
	if len(failures) != 1 || !strings.Contains(failures[0], "40 hex characters") {
		t.Fatalf("failures = %v, want '40 hex characters'", failures)
	}
}

func TestRequireHexStringRejectsInvalidHex(t *testing.T) {
	t.Parallel()
	var failures []string
	requireHexString(&failures, "commit", strings.Repeat("zz", 20), 20)
	if len(failures) != 1 || !strings.Contains(failures[0], "valid hex") {
		t.Fatalf("failures = %v, want 'valid hex'", failures)
	}
}

func TestValidateEvidenceRejectsEmptyGateName(t *testing.T) {
	t.Parallel()
	evidence := validEvidence()
	evidence.Gates = map[string]string{"": "passed"}
	failures := validateEvidence(evidence)
	if !containsFailure(failures, "empty gate name") {
		t.Fatalf("failures = %v, want empty gate name", failures)
	}
}

func TestValidateEvidenceRejectsEmptyGateStatus(t *testing.T) {
	t.Parallel()
	evidence := validEvidence()
	evidence.Gates = map[string]string{"ci": "  "}
	failures := validateEvidence(evidence)
	if !containsFailure(failures, "gates.ci is required") {
		t.Fatalf("failures = %v, want gates.ci is required", failures)
	}
}

func TestVerifyEvidenceRejectsNonexistentFile(t *testing.T) {
	t.Parallel()
	err := verifyEvidence("/nonexistent/path.json")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestVerifyEvidenceRejectsInvalidJSON(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "bad.json")
	os.WriteFile(path, []byte("not json"), 0o644)
	err := verifyEvidence(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestFlexibleStringUnmarshalNull(t *testing.T) {
	t.Parallel()
	var fs flexibleString
	if err := json.Unmarshal([]byte("null"), &fs); err != nil {
		t.Fatal(err)
	}
	if fs != "" {
		t.Fatalf("flexibleString(null) = %q, want empty", fs)
	}
}

func TestFlexibleStringUnmarshalEmpty(t *testing.T) {
	t.Parallel()
	var fs flexibleString
	if err := json.Unmarshal([]byte(`""`), &fs); err != nil {
		t.Fatal(err)
	}
	if fs != "" {
		t.Fatalf("flexibleString(\"\") = %q, want empty", fs)
	}
}

func TestFlexibleStringUnmarshalQuoted(t *testing.T) {
	t.Parallel()
	var fs flexibleString
	if err := json.Unmarshal([]byte(`"12345"`), &fs); err != nil {
		t.Fatal(err)
	}
	if fs != "12345" {
		t.Fatalf(`flexibleString("12345") = %q, want "12345"`, fs)
	}
}

func TestRunCLIHelpReturnsZero(t *testing.T) {
	t.Parallel()
	var stdout, stderr bytes.Buffer
	code := runCLI("downstreamadoption", []string{"-help"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("runCLI -help exit code = %d, want 0", code)
	}
}

type errorWriter struct{}

func (errorWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}
