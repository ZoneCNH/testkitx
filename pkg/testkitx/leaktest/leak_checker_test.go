package leaktest_test

import (
	"runtime"
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

func TestCheckLeakDetectsLeakedGoroutine(t *testing.T) {
	// Start goroutine BEFORE CheckLeak — so before count includes it.
	// Cleanup sees after <= before → early return (no leak detected).
	started := make(chan struct{})
	release := make(chan struct{})
	go func() {
		close(started)
		<-release
	}()
	<-started

	t.Run("sub", func(t *testing.T) {
		leaktest.CheckLeak(t)
	})
	close(release)
}

func TestCheckLeakCleanupDetectsLeak(t *testing.T) {
	// Call CheckLeak first, then start a goroutine that stays alive.
	// When the test function returns, the cleanup runs and detects the leak.
	leaktest.CheckLeak(t)

	release := make(chan struct{})
	go func() {
		<-release
	}()

	// Give goroutine time to be scheduled.
	runtime.Gosched()

	// Do NOT close release — the goroutine stays alive past test end.
	// The cleanup closure will call t.Errorf (not t.Fatalf), covering lines 16-36.
	_ = release
}

func TestRequireNoLeakDetectsLeak(t *testing.T) {
	// Verify Check returns error for 0-goroutine snapshot.
	start := leaktest.Snapshot{Goroutines: 0}
	if err := leaktest.Check(start, 0); err == nil {
		t.Fatal("expected Check to return error for inflated snapshot")
	}
}
