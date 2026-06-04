package main

import (
	"runtime"
	"time"
)

var checkNames = []string{
	"fmt",
	"vet",
	"lint",
	"unit_test",
	"race_test",
	"boundary",
	"secret_scan",
	"security",
	"contract",
	"integration",
}

var checkEnvNames = map[string]string{
	"fmt":         "FMT_STATUS",
	"vet":         "VET_STATUS",
	"lint":        "LINT_STATUS",
	"unit_test":   "UNIT_TEST_STATUS",
	"race_test":   "RACE_TEST_STATUS",
	"boundary":    "BOUNDARY_STATUS",
	"secret_scan": "SECRET_SCAN_STATUS",
	"security":    "SECURITY_STATUS",
	"contract":    "CONTRACT_STATUS",
	"integration": "INTEGRATION_STATUS",
}

type Manifest struct {
	Module           string            `json:"module"`
	Version          string            `json:"version"`
	Commit           string            `json:"commit"`
	TreeSHA          string            `json:"tree_sha"`
	SourceDigest     string            `json:"source_digest"`
	TrackedFileCount int               `json:"tracked_file_count"`
	GoVersion        string            `json:"go_version"`
	GeneratedAt      string            `json:"generated_at"`
	GeneratedBy      string            `json:"generated_by"`
	TreeState        string            `json:"tree_state"`
	Checks           map[string]string `json:"checks"`
	Contracts        []FileDigest      `json:"contracts"`
	Dependencies     []ModuleDigest    `json:"dependencies"`
	Tools            map[string]string `json:"tools"`
	Artifacts        []string          `json:"artifacts"`
	Notes            Notes             `json:"notes"`
}

type FileDigest struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

type ModuleDigest struct {
	Path    string         `json:"path"`
	Version string         `json:"version,omitempty"`
	Main    bool           `json:"main,omitempty"`
	Replace *ModuleReplace `json:"replace,omitempty"`
}

type ModuleReplace struct {
	Path    string `json:"path"`
	Version string `json:"version,omitempty"`
}

type Notes struct {
	BreakingChanges string   `json:"breaking_changes"`
	KnownRisks      []string `json:"known_risks"`
}

func buildManifest() (Manifest, error) {
	module, err := runTrimmed("go", "list", "-m")
	if err != nil {
		return Manifest{}, err
	}

	sourceDigest, trackedFileCount, err := sourceDigest()
	if err != nil {
		return Manifest{}, err
	}
	contracts, err := contractDigests()
	if err != nil {
		return Manifest{}, err
	}
	dependencies, err := moduleDigests()
	if err != nil {
		return Manifest{}, err
	}

	return Manifest{
		Module:           module,
		Version:          envDefault("VERSION", "v0.1.0"),
		Commit:           runTrimmedDefault("unknown", "git", "rev-parse", "HEAD"),
		TreeSHA:          runTrimmedDefault("unknown", "git", "rev-parse", "HEAD^{tree}"),
		SourceDigest:     sourceDigest,
		TrackedFileCount: trackedFileCount,
		GoVersion:        runtime.Version(),
		GeneratedAt:      time.Now().UTC().Format(time.RFC3339),
		GeneratedBy:      envDefault("GENERATED_BY", "scripts/generate_manifest.sh"),
		TreeState:        treeState(),
		Checks:           buildChecks(),
		Contracts:        contracts,
		Dependencies:     dependencies,
		Tools: map[string]string{
			"go":            firstLine(runTrimmedDefault(runtime.Version(), "go", "version")),
			"golangci-lint": toolVersion("golangci-lint", "--version"),
			"govulncheck":   toolVersion("govulncheck", "-version"),
		},
		Artifacts: []string{
			"release/manifest/latest.json",
			"release/manifest/latest.json.sha256",
		},
		Notes: Notes{
			BreakingChanges: "none",
			KnownRisks:      []string{},
		},
	}, nil
}

func buildChecks() map[string]string {
	defaultStatus := envDefault("CHECK_STATUS", "unknown")
	checks := make(map[string]string, len(checkNames))
	for _, name := range checkNames {
		checks[name] = envDefault(checkEnvNames[name], defaultStatus)
	}
	return checks
}
