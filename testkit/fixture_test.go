package testkit

import (
	"testing"
	"time"
)

func TestConfigBuildsValidFixture(t *testing.T) {
	t.Parallel()
	cfg := Config("fixture")
	if cfg.Name != "fixture" {
		t.Fatalf("unexpected name: %q", cfg.Name)
	}
	if cfg.Timeout != time.Second {
		t.Fatalf("unexpected timeout: %s", cfg.Timeout)
	}
	RequireNoError(t, cfg.Validate())
}

func TestRequireNoErrorAcceptsNil(t *testing.T) {
	t.Parallel()
	RequireNoError(t, nil)
}
