package leaktest_test

import (
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/leaktest"
)

func TestCheckPassesWithoutNewGoroutines(t *testing.T) {
	t.Parallel()
	start := leaktest.Capture()
	if err := leaktest.Check(start, 2); err != nil {
		t.Fatalf("unexpected leak failure: %v", err)
	}
}

func TestCheckDetectsLiveGoroutine(t *testing.T) {
	start := leaktest.Capture()
	started := make(chan struct{})
	release := make(chan struct{})
	go func() {
		close(started)
		<-release
	}()
	<-started
	if err := leaktest.Check(start, 0); err == nil {
		close(release)
		t.Fatal("expected goroutine leak to be detected")
	}
	close(release)
}

func TestCheckLeakRegistersCleanup(t *testing.T) {
	// CheckLeak registers a t.Cleanup that runs after the test.
	// This test verifies CheckLeak doesn't panic and the cleanup runs.
	leaktest.CheckLeak(t)
}

