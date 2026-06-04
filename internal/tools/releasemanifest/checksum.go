package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func writeManifest(path string, manifest Manifest) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(manifest); err != nil {
		return err
	}
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		return err
	}
	return writeManifestChecksum(path, buf.Bytes())
}

func manifestChecksumPath(path string) string {
	return path + ".sha256"
}

func writeManifestChecksum(path string, data []byte) error {
	sum := sha256.Sum256(data)
	content := fmt.Sprintf("%s  %s\n", hex.EncodeToString(sum[:]), filepath.Base(path))
	return os.WriteFile(manifestChecksumPath(path), []byte(content), 0o644)
}

func verifyManifestChecksum(path string, data []byte) error {
	sidecarPath := manifestChecksumPath(path)
	expected, err := readManifestChecksum(sidecarPath)
	if err != nil {
		return err
	}
	sum := sha256.Sum256(data)
	got := hex.EncodeToString(sum[:])
	if got != expected {
		return fmt.Errorf("%s does not match %s", sidecarPath, path)
	}
	return nil
}

func readManifestChecksum(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		return "", fmt.Errorf("%s is empty", path)
	}
	checksum := strings.TrimPrefix(strings.ToLower(fields[0]), "sha256:")
	if len(checksum) != sha256.Size*2 {
		return "", fmt.Errorf("%s has invalid sha256 digest", path)
	}
	if _, err := hex.DecodeString(checksum); err != nil {
		return "", fmt.Errorf("%s has invalid sha256 digest: %w", path, err)
	}
	return checksum, nil
}
