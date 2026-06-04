// Package golden provides opt-in golden file update helpers.
package golden

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

const UpdateEnv = "TESTKITX_UPDATE_GOLDEN"

type Evidence struct {
	Kind         string `json:"kind"`
	Path         string `json:"path"`
	Updated      bool   `json:"updated"`
	Matched      bool   `json:"matched"`
	ActualSHA256 string `json:"actual_sha256"`
}

func AssertBytes(t testing.TB, path string, actual []byte) Evidence {
	t.Helper()
	clean := filepath.Clean(path)
	actualHash := sha256Hex(actual)
	if os.Getenv(UpdateEnv) == "1" {
		if err := os.MkdirAll(filepath.Dir(clean), 0o755); err != nil {
			t.Fatalf("create golden dir: %v", err)
		}
		if err := os.WriteFile(clean, actual, 0o644); err != nil {
			t.Fatalf("update golden %s: %v", clean, err)
		}
		return Evidence{Kind: "golden_check", Path: clean, Updated: true, Matched: true, ActualSHA256: actualHash}
	}
	expected, err := os.ReadFile(clean)
	if err != nil {
		t.Fatalf("read golden %s: %v", clean, err)
	}
	if !bytes.Equal(expected, actual) {
		t.Fatalf("golden mismatch for %s\nexpected sha256: %s\nactual sha256:   %s", clean, sha256Hex(expected), actualHash)
	}
	return Evidence{Kind: "golden_check", Path: clean, Updated: false, Matched: true, ActualSHA256: actualHash}
}

func AssertJSON(t testing.TB, path string, value any) Evidence {
	t.Helper()
	encoded, err := canonicalJSON(value)
	if err != nil {
		t.Fatalf("canonicalize json: %v", err)
	}
	return AssertBytes(t, path, encoded)
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
