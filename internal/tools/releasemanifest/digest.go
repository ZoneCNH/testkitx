package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"sort"
	"strings"
)

func contractFiles() []string {
	return []string{
		"contracts/config.schema.json",
		"contracts/docker-toolchain.schema.json",
		"contracts/downstream-adoption-proof.schema.json",
		"contracts/error.schema.json",
		"contracts/health.schema.json",
		"contracts/metrics.md",
	}
}

func sourceDigest() (string, int, error) {
	raw, err := runRaw("git", "ls-files", "-z")
	if err != nil {
		return "", 0, err
	}
	parts := strings.Split(string(raw), "\x00")
	files := make([]string, 0, len(parts))
	for _, part := range parts {
		if part != "" {
			files = append(files, part)
		}
	}
	sort.Strings(files)

	digest := sha256.New()
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			return "", 0, err
		}
		fileSum := sha256.Sum256(data)
		digest.Write([]byte(path))
		digest.Write([]byte{0})
		digest.Write([]byte(hex.EncodeToString(fileSum[:])))
		digest.Write([]byte{0})
	}

	return "sha256:" + hex.EncodeToString(digest.Sum(nil)), len(files), nil
}

func contractDigests() ([]FileDigest, error) {
	digests := make([]FileDigest, 0, len(contractFiles()))
	for _, path := range contractFiles() {
		digest, err := fileDigest(path)
		if err != nil {
			return nil, err
		}
		digests = append(digests, digest)
	}
	return digests, nil
}

func fileDigest(path string) (FileDigest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return FileDigest{}, err
	}
	sum := sha256.Sum256(data)
	return FileDigest{
		Path:   path,
		SHA256: "sha256:" + hex.EncodeToString(sum[:]),
	}, nil
}

func moduleDigests() ([]ModuleDigest, error) {
	raw, err := runRaw("go", "list", "-m", "-json", "all")
	if err != nil {
		return nil, err
	}

	type goModule struct {
		Path    string
		Version string
		Main    bool
		Replace *struct {
			Path    string
			Version string
		}
	}

	decoder := json.NewDecoder(bytes.NewReader(raw))
	var modules []ModuleDigest
	for {
		var module goModule
		if err := decoder.Decode(&module); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		digest := ModuleDigest{
			Path:    module.Path,
			Version: module.Version,
			Main:    module.Main,
		}
		if module.Replace != nil {
			digest.Replace = &ModuleReplace{
				Path:    module.Replace.Path,
				Version: module.Replace.Version,
			}
		}
		modules = append(modules, digest)
	}
	return modules, nil
}
