package requirex

import "time"

func Eventually(t TestingT, timeout, interval time.Duration, fn func() bool) {
	t.Helper()
	if fn == nil {
		t.Fatalf("eventually predicate is required")
	}
	if timeout <= 0 {
		t.Fatalf("timeout must be positive, got %s", timeout)
	}
	if interval <= 0 {
		interval = time.Millisecond
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
