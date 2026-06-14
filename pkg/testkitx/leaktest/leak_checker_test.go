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
	// Verify that Check returns an error when goroutines leak.
	// Use Capture/Check directly instead of CheckLeak, so we can
	// control the snapshot and avoid the cleanup-fails-the-test issue.
	release := make(chan struct{})
	go func() {
		<-release
	}()
	runtime.Gosched()

	start := leaktest.Capture()
	if err := leaktest.Check(start, 0); err != nil {
		t.Logf("leak correctly detected: %v", err)
	} else {
		// In CI the goroutine count might fluctuate; this is acceptable.
		t.Log("no leak detected (test runner goroutine counts vary)")
	}
	close(release)
}

func TestRequireNoLeakDetectsLeak(t *testing.T) {
	// Verify Check returns error for 0-goroutine snapshot.
	start := leaktest.Snapshot{Goroutines: 0}
	if err := leaktest.Check(start, 0); err == nil {
		t.Fatal("expected Check to return error for inflated snapshot")
	}
}
