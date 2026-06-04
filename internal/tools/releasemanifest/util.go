package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func treeState() string {
	status, err := runTrimmed("git", "status", "--porcelain", "--untracked-files=all")
	if err != nil {
		return "unknown"
	}
	if status == "" {
		return "clean"
	}
	return "dirty"
}

func toolVersion(name string, args ...string) string {
	if _, err := exec.LookPath(name); err != nil {
		return "missing"
	}
	output, err := runTrimmed(name, args...)
	if err != nil {
		return "error: " + firstLine(err.Error())
	}
	return firstLine(output)
}

func runTrimmedDefault(fallback string, name string, args ...string) string {
	output, err := runTrimmed(name, args...)
	if err != nil {
		return fallback
	}
	return output
}

func runTrimmed(name string, args ...string) (string, error) {
	output, err := runRaw(name, args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func runRaw(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%s %s failed: %w: %s", name, strings.Join(args, " "), err, strings.TrimSpace(string(output)))
	}
	return output, nil
}

func envDefault(name string, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return fallback
}

func firstLine(value string) string {
	value = strings.TrimSpace(value)
	if idx := strings.IndexByte(value, '\n'); idx >= 0 {
		return value[:idx]
	}
	return value
}

func requireNonEmpty(failures *[]string, field string, value string) {
	if strings.TrimSpace(value) == "" {
		*failures = append(*failures, field+" is required")
	}
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
