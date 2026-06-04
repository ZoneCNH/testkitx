package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/ZoneCNH/testkitx/internal/cliutil"
)

const requiredCICommand = "GOWORK=off make ci"

type Evidence struct {
	Status            string            `json:"status"`
	Repository        string            `json:"repository"`
	Commit            string            `json:"commit"`
	TreeSHA           string            `json:"tree_sha"`
	WorkflowRunID     flexibleString    `json:"workflow_run_id"`
	WorkflowURL       string            `json:"workflow_url,omitempty"`
	ArtifactURL       string            `json:"artifact_url"`
	SHA256            string            `json:"sha256"`
	TestImports       int               `json:"test_imports"`
	ProductionImports int               `json:"production_imports"`
	Commands          []CommandEvidence `json:"commands"`
	Gates             map[string]string `json:"gates,omitempty"`
}

type CommandEvidence struct {
	Command string `json:"command"`
	Status  string `json:"status"`
}

type flexibleString string

func (value *flexibleString) UnmarshalJSON(data []byte) error {
	raw := strings.TrimSpace(string(data))
	if raw == "" || raw == "null" {
		*value = ""
		return nil
	}
	if strings.HasPrefix(raw, `"`) {
		var decoded string
		if err := json.Unmarshal(data, &decoded); err != nil {
			return err
		}
		*value = flexibleString(strings.TrimSpace(decoded))
		return nil
	}
	*value = flexibleString(raw)
	return nil
}

func main() {
	os.Exit(runCLI(os.Args[0], os.Args[1:], os.Stdout, os.Stderr))
}

func runCLI(name string, args []string, stdout io.Writer, stderr io.Writer) int {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(stderr)
	verify := flags.String("verify", "", "verify downstream adoption evidence JSON")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if strings.TrimSpace(*verify) == "" {
		return cliutil.PrintCLIError(stderr, errors.New("-verify is required"))
	}
	if err := verifyEvidence(*verify); err != nil {
		return cliutil.PrintCLIError(stderr, err)
	}
	return cliutil.PrintCLIStatus(stdout, "downstream adoption evidence verified: %s\n", *verify)
}

func verifyEvidence(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var evidence Evidence
	if err := json.Unmarshal(data, &evidence); err != nil {
		return err
	}
	if failures := validateEvidence(evidence); len(failures) > 0 {
		return errors.New("downstream adoption evidence verification failed:\n - " + strings.Join(failures, "\n - "))
	}
	return nil
}

func validateEvidence(evidence Evidence) []string {
	var failures []string
	requireEqual(&failures, "status", evidence.Status, "passed")
	requireNonEmpty(&failures, "repository", evidence.Repository)
	requireGitSHA(&failures, "commit", evidence.Commit)
	requireGitSHA(&failures, "tree_sha", evidence.TreeSHA)
	requireNumeric(&failures, "workflow_run_id", string(evidence.WorkflowRunID))
	requireHTTPURL(&failures, "artifact_url", evidence.ArtifactURL)
	if strings.TrimSpace(evidence.WorkflowURL) != "" {
		requireHTTPURL(&failures, "workflow_url", evidence.WorkflowURL)
	}
	requireSHA256Digest(&failures, "sha256", evidence.SHA256)

	if evidence.TestImports <= 0 {
		failures = append(failures, "test_imports must be greater than 0")
	}
	if evidence.ProductionImports != 0 {
		failures = append(failures, fmt.Sprintf("production_imports must be 0, got %d", evidence.ProductionImports))
	}
	if !hasPassedCommand(evidence.Commands, requiredCICommand) {
		failures = append(failures, `commands must include passed command "GOWORK=off make ci"`)
	}
	if len(evidence.Gates) == 0 {
		failures = append(failures, "gates must include at least one gate status")
	}
	for gate, status := range evidence.Gates {
		if strings.TrimSpace(gate) == "" {
			failures = append(failures, "gates must not contain an empty gate name")
		}
		if strings.TrimSpace(status) == "" {
			failures = append(failures, fmt.Sprintf("gates.%s is required", gate))
		}
	}
	return failures
}

func hasPassedCommand(commands []CommandEvidence, want string) bool {
	for _, command := range commands {
		if strings.TrimSpace(command.Command) == want && strings.TrimSpace(command.Status) == "passed" {
			return true
		}
	}
	return false
}

func requireNonEmpty(failures *[]string, field string, value string) {
	if strings.TrimSpace(value) == "" {
		*failures = append(*failures, field+" is required")
	}
}

func requireEqual(failures *[]string, field string, got string, want string) {
	got = strings.TrimSpace(got)
	if got != want {
		*failures = append(*failures, fmt.Sprintf("%s must be %q, got %q", field, want, got))
	}
}

func requireNumeric(failures *[]string, field string, value string) {
	value = strings.TrimSpace(value)
	if value == "" {
		*failures = append(*failures, field+" is required")
		return
	}
	for _, char := range value {
		if char < '0' || char > '9' {
			*failures = append(*failures, field+" must contain digits only")
			return
		}
	}
}

func requireHTTPURL(failures *[]string, field string, value string) {
	value = strings.TrimSpace(value)
	if value == "" {
		*failures = append(*failures, field+" is required")
		return
	}
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		*failures = append(*failures, field+" must be an absolute URL")
		return
	}
	if parsed.Scheme != "https" && parsed.Scheme != "http" {
		*failures = append(*failures, field+" must use http or https")
	}
}

func requireGitSHA(failures *[]string, field string, value string) {
	requireHexString(failures, field, strings.ToLower(strings.TrimSpace(value)), 20)
}

func requireSHA256Digest(failures *[]string, field string, value string) {
	value = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(value)), "sha256:")
	requireHexString(failures, field, value, sha256.Size)
}

func requireHexString(failures *[]string, field string, value string, byteSize int) {
	if len(value) != byteSize*2 {
		*failures = append(*failures, fmt.Sprintf("%s must be %d hex characters", field, byteSize*2))
		return
	}
	if _, err := hex.DecodeString(value); err != nil {
		*failures = append(*failures, field+" must be valid hex")
	}
}

