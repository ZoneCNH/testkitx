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

func TestCheckLeakWithIgnorePatterns(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	go func() {
		close(started)
		<-release
	}()
	<-started

	leaktest.CheckLeak(t, "leaktest_test.TestCheckLeakWithIgnorePatterns")
	close(release)
}

