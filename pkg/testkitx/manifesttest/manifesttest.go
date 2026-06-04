// Package manifesttest builds and validates release manifest fixtures.
package manifesttest

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
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
	return os.WriteFile(filepath.Clean(path), append(data, '\n'), 0o644)
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
