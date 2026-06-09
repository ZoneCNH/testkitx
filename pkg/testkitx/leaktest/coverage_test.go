package leaktest_test

import (
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/leaktest"
)

func TestRequireNoLeakPasses(t *testing.T) {
	t.Parallel()
	start := leaktest.Capture()
	leaktest.RequireNoLeak(t, start, 5)
}

func TestCaptureReturnsSnapshot(t *testing.T) {
	t.Parallel()
	s := leaktest.Capture()
	if s.Goroutines < 1 {
		t.Fatalf("expected at least 1 goroutine, got %d", s.Goroutines)
	}
}
