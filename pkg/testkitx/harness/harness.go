// Package harness runs commands and captures deterministic evidence.
package harness

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type Command struct {
	Dir     string
	Name    string
	Args    []string
	Env     map[string]string
	Timeout time.Duration
}

type Result struct {
	Kind         string   `json:"kind"`
	Command      []string `json:"command"`
	Dir          string   `json:"dir"`
	ExitCode     int      `json:"exit_code"`
	StdoutSHA256 string   `json:"stdout_sha256"`
	StderrSHA256 string   `json:"stderr_sha256"`
	EnvDigest    string   `json:"env_digest"`
	DurationMS   int64    `json:"duration_ms"`
	TimedOut     bool     `json:"timed_out"`
}

func Run(ctx context.Context, command Command) Result {
	if command.Timeout <= 0 {
		command.Timeout = 10 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, command.Timeout)
	defer cancel()
	started := time.Now()
	cmd := exec.CommandContext(ctx, command.Name, command.Args...)
	cmd.Dir = command.Dir
	for k, v := range command.Env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	stdout, stderr := &strings.Builder{}, &strings.Builder{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	exit := 0
	if err != nil {
		exit = 1
		if ee, ok := err.(*exec.ExitError); ok {
			exit = ee.ExitCode()
		}
	}
	return Result{Kind: "harness_command", Command: append([]string{command.Name}, command.Args...), Dir: command.Dir, ExitCode: exit, StdoutSHA256: digest(stdout.String()), StderrSHA256: digest(stderr.String()), EnvDigest: envDigest(command.Env), DurationMS: time.Since(started).Milliseconds(), TimedOut: ctx.Err() == context.DeadlineExceeded}
}

func envDigest(env map[string]string) string {
	parts := make([]string, 0, len(env))
	for k, v := range env {
		parts = append(parts, k+"="+v)
	}
	sort.Strings(parts)
	return digest(strings.Join(parts, "\n"))
}

func digest(s string) string { sum := sha256.Sum256([]byte(s)); return hex.EncodeToString(sum[:]) }
