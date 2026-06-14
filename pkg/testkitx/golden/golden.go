// Package golden provides opt-in golden file update helpers.
package golden

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

const UpdateEnv = "TESTKITX_UPDATE_GOLDEN"

// GoldenUpdate reports whether golden files should be updated.
// It checks the GOLDEN_UPDATE environment variable per SPEC §11 (FR-008).
// It also respects TESTKITX_UPDATE_GOLDEN for backward compatibility.
func GoldenUpdate() bool {
	return os.Getenv("GOLDEN_UPDATE") == "1" || os.Getenv(UpdateEnv) == "1"
}

type Evidence struct {
	Kind         string `json:"kind"`
	Path         string `json:"path"`
	Updated      bool   `json:"updated"`
	Matched      bool   `json:"matched"`
	ActualSHA256 string `json:"actual_sha256"`
}

func AssertBytes(t testing.TB, path string, actual []byte) Evidence {
	t.Helper()
	ev, err := CheckBytes(path, actual)
	if err != nil {
		t.Fatalf("%v", err)
	}
	return ev
}

// CheckBytes compares actual against the golden file at path.
// Returns an error instead of failing the test.
func CheckBytes(path string, actual []byte) (Evidence, error) {
	clean := filepath.Clean(path)
	actualHash := sha256Hex(actual)
	if os.Getenv(UpdateEnv) == "1" {
		if err := os.MkdirAll(filepath.Dir(clean), 0o755); err != nil {
			return Evidence{}, fmt.Errorf("create golden dir: %w", err)
		}
		if err := os.WriteFile(clean, actual, 0o644); err != nil {
			return Evidence{}, fmt.Errorf("update golden %s: %w", clean, err)
		}
		return Evidence{Kind: "golden_check", Path: clean, Updated: true, Matched: true, ActualSHA256: actualHash}, nil
	}
	expected, err := os.ReadFile(clean)
	if err != nil {
		return Evidence{}, fmt.Errorf("read golden %s: %w", clean, err)
	}
	if !bytes.Equal(expected, actual) {
		return Evidence{}, fmt.Errorf("golden mismatch for %s\nexpected sha256: %s\nactual sha256:   %s", clean, sha256Hex(expected), actualHash)
	}
	return Evidence{Kind: "golden_check", Path: clean, Updated: false, Matched: true, ActualSHA256: actualHash}, nil
}

func AssertJSON(t testing.TB, path string, value any) Evidence {
	t.Helper()
	ev, err := CheckJSON(path, value)
	if err != nil {
		t.Fatalf("%v", err)
	}
	return ev
}

// CheckJSON compares value against the golden file at path using canonical JSON.
// Returns an error instead of failing the test.
func CheckJSON(path string, value any) (Evidence, error) {
	encoded, err := canonicalJSON(value)
	if err != nil {
		return Evidence{}, fmt.Errorf("canonicalize json: %w", err)
	}
	return CheckBytes(path, encoded)
}

func canonicalJSON(value any) ([]byte, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var decoded any
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, err
	}
	return json.MarshalIndent(decoded, "", "  ")
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
