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
	if err := os.MkdirAll(filepath.Dir(filepath.Clean(path)), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	cleanPath := filepath.Clean(path)
	if err := os.WriteFile(cleanPath, append(data, '\n'), 0o644); err != nil {
		return err
	}
	return WriteChecksum(cleanPath)
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

func ChecksumPath(path string) string {
	return filepath.Clean(path) + ".sha256"
}

func SHA256(path string) (string, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func WriteChecksum(path string) error {
	digest, err := SHA256(path)
	if err != nil {
		return err
	}
	rawDigest := strings.TrimPrefix(digest, "sha256:")
	content := fmt.Sprintf("%s  %s\n", rawDigest, filepath.Base(path))
	return os.WriteFile(ChecksumPath(path), []byte(content), 0o644)
}

func VerifyChecksum(path string) error {
	digest, err := SHA256(path)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(ChecksumPath(path))
	if err != nil {
		return err
	}
	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		return errors.New("manifest checksum is empty")
	}
	want := strings.TrimPrefix(strings.ToLower(fields[0]), "sha256:")
	got := strings.TrimPrefix(digest, "sha256:")
	if got != want {
		return fmt.Errorf("manifest checksum mismatch: got %s, want %s", got, want)
	}
	return nil
}

func AssertChecksum(t testing.TB, path string) {
	t.Helper()
	if err := VerifyChecksum(path); err != nil {
		t.Fatalf("manifest checksum verification failed: %v", err)
	}
}
