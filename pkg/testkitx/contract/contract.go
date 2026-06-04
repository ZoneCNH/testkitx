// Package contract provides contract hash assertions with machine-readable evidence.
package contract

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type Evidence struct {
	Kind         string            `json:"kind"`
	ContractID   string            `json:"contract_id"`
	ContractPath string            `json:"contract_path"`
	SHA256       string            `json:"sha256"`
	Matched      bool              `json:"matched"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

func AssertHash(t testing.TB, contractID, path, expectedSHA256 string, metadata map[string]string) Evidence {
	t.Helper()
	actual, err := FileSHA256(path)
	if err != nil {
		t.Fatalf("hash contract %s: %v", path, err)
	}
	matched := actual == expectedSHA256
	if !matched {
		t.Fatalf("contract %s hash mismatch: got %s want %s", contractID, actual, expectedSHA256)
	}
	return Evidence{Kind: "contract_check", ContractID: contractID, ContractPath: filepath.Clean(path), SHA256: actual, Matched: true, Metadata: metadata}
}

func WriteEvidence(path string, evidence Evidence) error {
	if err := os.MkdirAll(filepath.Dir(filepath.Clean(path)), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(evidence, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Clean(path), append(data, '\n'), 0o644)
}

func FileSHA256(path string) (string, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}
