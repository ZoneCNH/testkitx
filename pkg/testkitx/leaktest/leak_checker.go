package leaktest

import (
	"runtime"
	"strings"
	"testing"
)

// CheckLeak captures the goroutine count at call time and registers a t.Cleanup
// that compares against the snapshot. If the count grew and the new goroutines
// do not match any ignore pattern, the test is failed with a full stack dump.
func CheckLeak(t *testing.T, ignorePatterns ...string) {
	t.Helper()
	before := Capture()

	t.Cleanup(func() {
		// Force GC to flush finalizers.
		for i := 0; i < 20; i++ {
			runtime.GC()
		}

		after := runtime.NumGoroutine()
		if after <= before.Goroutines {
			return
		}

		buf := make([]byte, 1<<20)
		n := runtime.Stack(buf, true)
		stack := string(buf[:n])

		unexpected := countUnexpectedGoroutines(stack, ignorePatterns)
		if unexpected > 0 {
			t.Errorf("goroutine leak: before=%d after=%d unexpected=%d\n%s",
				before.Goroutines, after, unexpected, stack)
		}
	})
}

// IgnoreGoroutines is a convenience that returns the patterns slice unchanged.
// It documents intent at the call site:
//
//	leaktest.CheckLeak(t, leaktest.IgnoreGoroutines("runtime.", "testing.")...)
func IgnoreGoroutines(patterns ...string) []string {
	return patterns
}

// countUnexpectedGoroutines splits a full goroutine dump into blocks and counts
// those that do not match any ignore pattern.
func countUnexpectedGoroutines(stack string, ignorePatterns []string) int {
	blocks := strings.Split(stack, "\n\n")
	count := 0
	for _, block := range blocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}
		matched := false
		for _, pat := range ignorePatterns {
			if strings.Contains(block, pat) {
				matched = true
				break
			}
		}
		if !matched {
			count++
		}
	}
	return count
}
