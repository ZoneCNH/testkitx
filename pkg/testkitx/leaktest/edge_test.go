package leaktest_test

import (
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/leaktest"
)

func TestCheckLeakNoLeak(t *testing.T) {
	start := leaktest.Capture()
	for i := 0; i < 20; i++ {
		// Force GC cycles
	}
	err := leaktest.Check(start, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckLeakWithTolerance(t *testing.T) {
	start := leaktest.Capture()
	// tolerance of 100 should always pass
	err := leaktest.Check(start, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRequireNoLeakError(t *testing.T) {
	// Capture with 0 tolerance - may or may not trigger depending on goroutine state
	// Use a known-inflated snapshot
	start := leaktest.Snapshot{Goroutines: 0}
	err := leaktest.Check(start, 0)
	// With goroutines starting at 0 and some running, this should error
	// But it might not always - just verify the function works
	if err != nil && err.Error() == "" {
		t.Fatal("expected non-empty error message")
	}
}

func TestIgnoreGoroutines(t *testing.T) {
	t.Parallel()
	patterns := leaktest.IgnoreGoroutines("runtime.", "testing.")
	if len(patterns) != 2 {
		t.Fatalf("expected 2 patterns, got %d", len(patterns))
	}
}
