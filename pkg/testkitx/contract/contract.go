// Package contract provides contract hash assertions with machine-readable evidence.
package contract

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

func (e Evidence) Validate() error {
	if strings.TrimSpace(e.Kind) == "" {
		return fmt.Errorf("kind is required")
	}
	if e.Kind != "contract_check" {
		return fmt.Errorf("kind must be contract_check")
	}
	if strings.TrimSpace(e.ContractID) == "" {
		return fmt.Errorf("contract_id is required")
	}
	if strings.TrimSpace(e.ContractPath) == "" {
		return fmt.Errorf("contract_path is required")
	}
	if strings.TrimSpace(e.SHA256) == "" {
		return fmt.Errorf("sha256 is required")
	}
	if len(e.SHA256) != sha256.Size*2 {
		return fmt.Errorf("sha256 is invalid")
	}
	if _, err := hex.DecodeString(e.SHA256); err != nil {
		return fmt.Errorf("sha256 is invalid: %w", err)
	}
	if !e.Matched {
		return fmt.Errorf("matched must be true")
	}
	return nil
}

func AssertHash(t testing.TB, contractID, path, expectedSHA256 string, metadata map[string]string) Evidence {
	t.Helper()
	if strings.TrimSpace(contractID) == "" {
		t.Fatalf("contract id is required")
	}
	actual, err := FileSHA256(path)
	if err != nil {
		t.Fatalf("hash contract %s: %v", path, err)
	}
	matched := actual == expectedSHA256
	if !matched {
		t.Fatalf("contract %s hash mismatch: got %s want %s", contractID, actual, expectedSHA256)
	}
	return Evidence{Kind: "contract_check", ContractID: contractID, ContractPath: filepath.Clean(path), SHA256: actual, Matched: true, Metadata: copyMetadata(metadata)}
}

func WriteEvidence(path string, evidence Evidence) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("evidence path is required")
	}
	if err := evidence.Validate(); err != nil {
		return err
	}
	cleanPath := filepath.Clean(path)
	if err := os.MkdirAll(filepath.Dir(cleanPath), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(evidence, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cleanPath, append(data, '\n'), 0o644)
}

func FileSHA256(path string) (string, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func copyMetadata(metadata map[string]string) map[string]string {
	if len(metadata) == 0 {
		return nil
	}
	copied := make(map[string]string, len(metadata))
	for key, value := range metadata {
		copied[key] = value
	}
	return copied
}
