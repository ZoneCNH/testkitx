package testkitx

import (
	"testing"
	"time"
)

// Eventually polls fn until it returns true or timeout expires, then fails t.
// Per SPEC FR-007.
//
// When fn returns true within timeout, the test continues.
// When fn times out, t is failed with a clear diagnostic.
// timeout <= 0 means check exactly once.
func Eventually(t *testing.T, fn func() bool, timeout, interval time.Duration) {
	t.Helper()
	if fn == nil {
		t.Fatalf("Eventually: predicate must not be nil")
	}
	if timeout <= 0 {
		if fn() {
			return
		}
		t.Fatalf("condition was not satisfied (checked once, timeout=%s)", timeout)
		return
	}
	if interval <= 0 {
		interval = 10 * time.Millisecond
	}
	deadline := time.Now().Add(timeout)
	for {
		if fn() {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("condition was not satisfied within %s", timeout)
		}
		sleep := interval
		if remaining := time.Until(deadline); remaining < sleep {
			sleep = remaining
		}
		if sleep > 0 {
			time.Sleep(sleep)
		}
	}
}
