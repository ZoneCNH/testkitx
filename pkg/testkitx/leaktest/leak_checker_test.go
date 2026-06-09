package leaktest_test

import (
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/leaktest"
)

func TestCheckLeakPassesWithoutNewGoroutines(t *testing.T) {
	leaktest.CheckLeak(t)
}

func TestIgnoreGoroutinesReturnsPatterns(t *testing.T) {
	t.Parallel()
	pats := leaktest.IgnoreGoroutines("runtime.", "testing.")
	if len(pats) != 2 {
		t.Fatalf("expected 2 patterns, got %d", len(pats))
	}
}
