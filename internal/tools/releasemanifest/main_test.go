package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ZoneCNH/testkitx/internal/cliutil"
)

func TestRunCLIGeneratesManifestToOut(t *testing.T) {
	t.Setenv("GOWORK", "off")
	t.Setenv("VERSION", "v1.2.3-cli")
	t.Setenv("GENERATED_BY", "releasemanifest-cli-test")
	t.Setenv("CHECK_STATUS", "passed")
	chdir(t, releaseManifestFixtureRepo(t))

	outPath := filepath.Join(t.TempDir(), "custom", "latest.json")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI("releasemanifest", []string{"-out", outPath}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("runCLI generate exit code = %d, want 0; stderr: %s", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	if want := "generated " + outPath; !strings.Contains(stdout.String(), want) {
		t.Fatalf("stdout = %q, want substring %q", stdout.String(), want)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(data) {
		t.Fatalf("generated manifest is invalid JSON: %s", data)
	}
	assertManifestChecksum(t, outPath)

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.Module != "example.com/releasefixture" {
		t.Fatalf("module = %q, want fixture module", manifest.Module)
	}
	if manifest.Version != "v1.2.3-cli" {
		t.Fatalf("version = %q, want v1.2.3-cli", manifest.Version)
	}
	if manifest.GeneratedBy != "releasemanifest-cli-test" {
		t.Fatalf("generated_by = %q, want releasemanifest-cli-test", manifest.GeneratedBy)
	}
	for _, name := range checkNames() {
		if manifest.Checks[name] != "passed" {
			t.Fatalf("checks[%q] = %q, want passed", name, manifest.Checks[name])
		}
	}
}

func TestRunCLIGenerateReportsBuildManifestFailure(t *testing.T) {
	t.Setenv("GOWORK", "off")
	repo := t.TempDir()
	runTestCommand(t, repo, "git", "init")
	if err := os.WriteFile(filepath.Join(repo, "go.mod"), []byte("module example.com/brokenmanifest\n\ngo 1.23\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runTestCommand(t, repo, "git", "add", ".")
	chdir(t, repo)

	outPath := filepath.Join(t.TempDir(), "latest.json")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI("releasemanifest", []string{"-out", outPath}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("runCLI generate exit code = %d, want 1; stdout: %s; stderr: %s", code, stdout.String(), stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	message := stderr.String()
	for _, want := range []string{"ERROR:", "contracts/config.schema.json"} {
		if !strings.Contains(message, want) {
			t.Fatalf("stderr = %q, want substring %q", message, want)
		}
	}
	if _, err := os.Stat(outPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("generated manifest exists after failed build: %v", err)
	}
	if _, err := os.Stat(manifestChecksumPath(outPath)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("generated manifest checksum exists after failed build: %v", err)
	}
}

func TestRunCLIGenerateReportsWriteManifestFailure(t *testing.T) {
	t.Setenv("GOWORK", "off")
	t.Setenv("CHECK_STATUS", "passed")
	chdir(t, releaseManifestFixtureRepo(t))

	outPath := t.TempDir()
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI("releasemanifest", []string{"-out", outPath}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("runCLI generate exit code = %d, want 1; stdout: %s; stderr: %s", code, stdout.String(), stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if want := "ERROR:"; !strings.Contains(stderr.String(), want) {
		t.Fatalf("stderr = %q, want substring %q", stderr.String(), want)
	}
}

func TestRunCLIVerifiesManifestWithRequirePassed(t *testing.T) {
	t.Setenv("GOWORK", "off")
	t.Setenv("VERSION", "v1.2.3")
	t.Setenv("CHECK_STATUS", "passed")
	chdir(t, releaseManifestFixtureRepo(t))

	outPath := filepath.Join(t.TempDir(), "latest.json")
	var generateStdout bytes.Buffer
	var generateStderr bytes.Buffer
	if code := runCLI("releasemanifest", []string{"-out", outPath}, &generateStdout, &generateStderr); code != 0 {
		t.Fatalf("runCLI generate exit code = %d, want 0; stderr: %s", code, generateStderr.String())
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runCLI("releasemanifest", []string{"-verify", outPath, "-require-passed", "-expect-version", "v1.2.3"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("runCLI verify exit code = %d, want 0; stderr: %s", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	if want := "release evidence verified: " + outPath; !strings.Contains(stdout.String(), want) {
		t.Fatalf("stdout = %q, want substring %q", stdout.String(), want)
	}
}

func TestRunCLIVerifyRejectsExpectedVersionMismatch(t *testing.T) {
	t.Setenv("GOWORK", "off")
	t.Setenv("VERSION", "v1.2.3")
	t.Setenv("CHECK_STATUS", "passed")
	chdir(t, releaseManifestFixtureRepo(t))

	outPath := filepath.Join(t.TempDir(), "latest.json")
	var generateStdout bytes.Buffer
	var generateStderr bytes.Buffer
	if code := runCLI("releasemanifest", []string{"-out", outPath}, &generateStdout, &generateStderr); code != 0 {
		t.Fatalf("runCLI generate exit code = %d, want 0; stderr: %s", code, generateStderr.String())
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runCLI("releasemanifest", []string{"-verify", outPath, "-expect-version", "v9.9.9"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("runCLI verify exit code = %d, want 1; stdout: %s; stderr: %s", code, stdout.String(), stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if want := `version mismatch: got "v1.2.3", want "v9.9.9"`; !strings.Contains(stderr.String(), want) {
		t.Fatalf("stderr = %q, want substring %q", stderr.String(), want)
	}
}

func TestRunCLIVerifyRejectsChecksumMismatch(t *testing.T) {
	t.Setenv("GOWORK", "off")
	t.Setenv("VERSION", "v1.2.3")
	t.Setenv("CHECK_STATUS", "passed")
	chdir(t, releaseManifestFixtureRepo(t))

	outPath := filepath.Join(t.TempDir(), "latest.json")
	var generateStdout bytes.Buffer
	var generateStderr bytes.Buffer
	if code := runCLI("releasemanifest", []string{"-out", outPath}, &generateStdout, &generateStderr); code != 0 {
		t.Fatalf("runCLI generate exit code = %d, want 0; stderr: %s", code, generateStderr.String())
	}
	if err := os.WriteFile(manifestChecksumPath(outPath), []byte(strings.Repeat("0", sha256.Size*2)+"  latest.json\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runCLI("releasemanifest", []string{"-verify", outPath}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("runCLI verify exit code = %d, want 1; stdout: %s; stderr: %s", code, stdout.String(), stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if want := manifestChecksumPath(outPath) + " does not match " + outPath; !strings.Contains(stderr.String(), want) {
		t.Fatalf("stderr = %q, want substring %q", stderr.String(), want)
	}
}

func TestRunCLIVerifyReportsDrift(t *testing.T) {
	t.Setenv("GOWORK", "off")
	t.Setenv("CHECK_STATUS", "passed")
	chdir(t, releaseManifestFixtureRepo(t))

	outPath := filepath.Join(t.TempDir(), "latest.json")
	var generateStdout bytes.Buffer
	var generateStderr bytes.Buffer
	if code := runCLI("releasemanifest", []string{"-out", outPath}, &generateStdout, &generateStderr); code != 0 {
		t.Fatalf("runCLI generate exit code = %d, want 0; stderr: %s", code, generateStderr.String())
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	manifest.SourceDigest = "sha256:stale"
	manifest.Checks["lint"] = "failed"
	if err := writeManifest(outPath, manifest); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runCLI("releasemanifest", []string{"-verify", outPath, "-require-passed"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("runCLI verify exit code = %d, want 1; stdout: %s; stderr: %s", code, stdout.String(), stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	message := stderr.String()
	for _, want := range []string{
		"ERROR: release evidence verification failed",
		"source_digest does not match current tracked file contents",
		`checks.lint must be passed, got "failed"`,
	} {
		if !strings.Contains(message, want) {
			t.Fatalf("stderr = %q, want substring %q", message, want)
		}
	}
}

func TestRunCLIVerifyRequiresCleanTree(t *testing.T) {
	t.Setenv("GOWORK", "off")
	t.Setenv("CHECK_STATUS", "passed")
	chdir(t, releaseManifestFixtureRepo(t))

	outPath := filepath.Join(t.TempDir(), "latest.json")
	var generateStdout bytes.Buffer
	var generateStderr bytes.Buffer
	if code := runCLI("releasemanifest", []string{"-out", outPath}, &generateStdout, &generateStderr); code != 0 {
		t.Fatalf("runCLI generate exit code = %d, want 0; stderr: %s", code, generateStderr.String())
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	manifest.TreeState = "dirty"
	if err := writeManifest(outPath, manifest); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runCLI("releasemanifest", []string{"-verify", outPath, "-require-clean"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("runCLI verify exit code = %d, want 1; stdout: %s; stderr: %s", code, stdout.String(), stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if want := `tree_state must be clean, got "dirty"`; !strings.Contains(stderr.String(), want) {
		t.Fatalf("stderr = %q, want substring %q", stderr.String(), want)
	}
}

func TestRunCLIHelpReturnsSuccess(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI("releasemanifest", []string{"-h"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("runCLI help exit code = %d, want 0; stderr: %s", code, stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if want := "Usage of releasemanifest"; !strings.Contains(stderr.String(), want) {
		t.Fatalf("stderr = %q, want substring %q", stderr.String(), want)
	}
}

func TestRunCLIRejectsUnknownFlag(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI("releasemanifest", []string{"-unknown"}, &stdout, &stderr)
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
	if code := cliutil.PrintCLIStatus(errorWriter{}, "ok\n"); code != 1 {
		t.Fatalf("printCLIStatus exit code = %d, want 1", code)
	}
	if code := cliutil.PrintCLIError(errorWriter{}, errors.New("boom")); code != 1 {
		t.Fatalf("printCLIError exit code = %d, want 1", code)
	}
}

// Shared test helpers

func repoRoot(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			if _, err := os.Stat(filepath.Join(dir, "contracts")); err == nil {
				return dir
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find repository root")
		}
		dir = parent
	}
}

func chdir(t *testing.T, dir string) {
	t.Helper()

	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(previous); err != nil {
			t.Fatalf("restore working directory: %v", err)
		}
	})
}

func runTestCommand(t *testing.T, dir string, name string, args ...string) {
	t.Helper()

	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %s failed: %v: %s", name, strings.Join(args, " "), err, strings.TrimSpace(string(output)))
	}
}

func releaseManifestFixtureRepo(t *testing.T) string {
	t.Helper()

	repo := t.TempDir()
	runTestCommand(t, repo, "git", "init")
	if err := os.WriteFile(filepath.Join(repo, "go.mod"), []byte("module example.com/releasefixture\n\ngo 1.23\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, path := range contractFiles() {
		fullPath := filepath.Join(repo, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatal(err)
		}
		content := "{}\n"
		if strings.HasSuffix(path, ".md") {
			content = "# Fixture Contract\n"
		}
		if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	runTestCommand(t, repo, "git", "add", ".")
	return repo
}

func assertManifestChecksum(t *testing.T, path string) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	sidecar, err := os.ReadFile(manifestChecksumPath(path))
	if err != nil {
		t.Fatal(err)
	}
	fields := strings.Fields(string(sidecar))
	if len(fields) < 1 {
		t.Fatalf("checksum sidecar is empty: %q", sidecar)
	}
	sum := sha256.Sum256(data)
	want := hex.EncodeToString(sum[:])
	if got := fields[0]; got != want {
		t.Fatalf("checksum = %q, want %q", got, want)
	}
	if len(fields) > 1 && fields[1] != filepath.Base(path) {
		t.Fatalf("checksum filename = %q, want %q", fields[1], filepath.Base(path))
	}
}

type errorWriter struct{}

func (errorWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}
