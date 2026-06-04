package harness_test

import (
	"context"
	"testing"
	"time"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/harness"
)

func TestRunCapturesCommandEvidence(t *testing.T) {
	t.Parallel()
	result := harness.Run(context.Background(), harness.Command{
		Name: "/bin/sh",
		Args: []string{"-c", "printf ok; printf err >&2; exit 7"},
		Env:  map[string]string{"TESTKITX_CASE": "harness"},
	})
	if result.Kind != "harness_command" || result.ExitCode != 7 || result.TimedOut {
		t.Fatalf("unexpected result: %+v", result)
	}
	if result.StdoutSHA256 == "" || result.StderrSHA256 == "" || result.EnvDigest == "" {
		t.Fatalf("missing digest evidence: %+v", result)
	}
	if len(result.Command) != 3 || result.Command[0] != "/bin/sh" {
		t.Fatalf("unexpected command evidence: %+v", result.Command)
	}
}

func TestRunMarksTimeout(t *testing.T) {
	t.Parallel()
	result := harness.Run(context.Background(), harness.Command{
		Name:    "/bin/sh",
		Args:    []string{"-c", "sleep 1"},
		Timeout: 10 * time.Millisecond,
	})
	if !result.TimedOut || result.ExitCode == 0 {
		t.Fatalf("expected timeout failure evidence, got %+v", result)
	}
}
