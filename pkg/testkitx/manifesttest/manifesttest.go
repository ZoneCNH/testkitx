// Package manifesttest builds and validates release manifest fixtures.
package manifesttest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type Manifest struct {
	Kind     string            `json:"kind"`
	Module   string            `json:"module"`
	Commit   string            `json:"commit"`
	Gates    map[string]string `json:"gates"`
	Evidence []string          `json:"evidence"`
}

func New(module, commit string) Manifest {
	return Manifest{Kind: "manifest_fixture_check", Module: module, Commit: commit, Gates: map[string]string{}, Evidence: []string{}}
}
func (m Manifest) Validate() error {
	if m.Kind != "manifest_fixture_check" || m.Module == "" || m.Commit == "" {
		return errors.New("manifest missing required fields")
	}
	return nil
}
func Write(path string, m Manifest) error {
	if err := m.Validate(); err != nil {
		return err
	}
	cleanPath := filepath.Clean(path)
	if err := os.MkdirAll(filepath.Dir(cleanPath), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(cleanPath, append(data, '\n'), 0o644); err != nil {
		return err
	}
	return WriteChecksum(cleanPath, ChecksumPath(cleanPath))
}
func Read(path string) (Manifest, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return Manifest{}, err
	}
	var m Manifest
	err = json.Unmarshal(data, &m)
	return m, err
}

func ChecksumPath(manifestPath string) string {
	return filepath.Clean(manifestPath) + ".sha256"
}

func SHA256(path string) (string, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func WriteChecksum(manifestPath string, checksumPath string) error {
	sum, err := SHA256(manifestPath)
	if err != nil {
		return err
	}
	if strings.TrimSpace(checksumPath) == "" {
		checksumPath = ChecksumPath(manifestPath)
	}
	content := fmt.Sprintf("%s  %s\n", sum, filepath.Base(filepath.Clean(manifestPath)))
	return os.WriteFile(filepath.Clean(checksumPath), []byte(content), 0o644)
}

func VerifyChecksum(manifestPath string, checksumPath string) error {
	got, err := SHA256(manifestPath)
	if err != nil {
		return err
	}
	if strings.TrimSpace(checksumPath) == "" {
		checksumPath = ChecksumPath(manifestPath)
	}
	data, err := os.ReadFile(filepath.Clean(checksumPath))
	if err != nil {
		return err
	}
	want := strings.Fields(string(data))
	if len(want) == 0 {
		return fmt.Errorf("empty checksum file: %s", checksumPath)
	}
	expected := strings.TrimPrefix(want[0], "sha256:")
	if len(expected) != sha256.Size*2 {
		return fmt.Errorf("invalid sha256 in %s", checksumPath)
	}
	if _, err := hex.DecodeString(expected); err != nil {
		return fmt.Errorf("invalid sha256 in %s: %w", checksumPath, err)
	}
	if expected != got {
		return fmt.Errorf("checksum mismatch for %s: got %s, want %s", manifestPath, got, expected)
	}
	return nil
}

func AssertManifestValid(t testing.TB, path string) {
	t.Helper()
	manifest, err := Read(path)
	if err != nil {
		t.Fatalf("read manifest fixture: %v", err)
	}
	if err := manifest.Validate(); err != nil {
		t.Fatalf("validate manifest fixture: %v", err)
	}
}

func AssertChecksum(t testing.TB, manifestPath string, checksumPath string) {
	t.Helper()
	if err := VerifyChecksum(manifestPath, checksumPath); err != nil {
		t.Fatalf("verify manifest checksum: %v", err)
	}
}
