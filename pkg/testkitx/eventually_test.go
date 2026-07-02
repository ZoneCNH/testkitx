package testkitx

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

// fakeTB implements testing.TB enough to assert Eventually's fail paths.
//
// testing.TB has unexported methods, so we satisfy the interface by embedding
// the interface value itself. The methods we actually exercise (Helper,
// Fatalf, Logf, Failed) are overridden as concrete methods below; any other
// method called through the embedded interface would panic, which surfaces a
// clear test failure if Eventually ever starts using them.
type fakeTB struct {
	testing.TB

	mu      sync.Mutex
	failed  bool
	logBuf  bytes.Buffer
}

func (f *fakeTB) Helper() {}

func (f *fakeTB) Fatalf(format string, args ...any) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.failed = true
	// emulate testing.T.Fatalf: record message then runtime.Goexit so the
	// surrounding goroutine stops the same way a real Fatalf would.
	f.logBuf.WriteString("[FATAL] " + fmt.Sprintf(format, args...) + "\n")
	runtime.Goexit()
}

func (f *fakeTB) Fatal(args ...any) {
	f.Fatalf("%s", fmt.Sprint(args...))
}

func (f *fakeTB) Logf(format string, args ...any) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.logBuf.WriteString("[LOG] " + fmt.Sprintf(format, args...) + "\n")
}

func (f *fakeTB) Log(args ...any) {
	f.Logf("%s", fmt.Sprint(args...))
}

func (f *fakeTB) Failed() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.failed
}

// Buffer returns a copy of the captured log output. Renamed to avoid clashing
// with testing.TB.Output() (which returns io.Writer).
func (f *fakeTB) Buffer() []byte {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]byte, f.logBuf.Len())
	copy(out, f.logBuf.Bytes())
	return out
}

// runEventuallyInGoroutine invokes Eventually on a fakeTB inside a separate
// goroutine and recovers the runtime.Goexit panic that Fatalf raises. It
// returns the fakeTB so the caller can inspect Failed() and captured output.
//
// Why a goroutine: testing.T.Fatalf calls runtime.Goexit. If we ran Eventually
// directly on the test's own *testing.T, the test goroutine would terminate and
// we could not assert afterwards. Running on a fakeTB inside a child goroutine
// keeps the test goroutine alive while still exercising the real fail path.
func runEventuallyInGoroutine(fn func() bool, timeout, interval time.Duration) *fakeTB {
	ftb := &fakeTB{}
	done := make(chan struct{})
	go func() {
		defer func() {
			recover() // swallow runtime.Goexit from Fatalf
			close(done)
		}()
		Eventually(ftb, fn, timeout, interval)
	}()
	<-done
	return ftb
}

// ===== TestEventually_* — closes FR-007 / AC-TKX-007 / TC-007 =====

// TestEventually_SucceedsWithinTimeout: fn returns true before timeout ->
// Eventually returns normally, t is NOT failed.
func TestEventually_SucceedsWithinTimeout(t *testing.T) {
	calls := 0
	Eventually(t, func() bool {
		calls++
		return calls >= 2
	}, 200*time.Millisecond, 5*time.Millisecond)

	if calls < 2 {
		t.Errorf("predicate should have been polled at least twice, got %d", calls)
	}
	if t.Failed() {
		t.Error("test should not be marked failed when predicate succeeds")
	}
}

// TestEventually_TimeoutFails: predicate never true -> fakeTB failed with a
// diagnostic that mentions the timeout window.
func TestEventually_TimeoutFails(t *testing.T) {
	ftb := runEventuallyInGoroutine(
		func() bool { return false },
		30*time.Millisecond,
		5*time.Millisecond,
	)

	if !ftb.Failed() {
		t.Fatal("Eventually should fail when predicate never becomes true")
	}
	out := string(ftb.Buffer())
	if !strings.Contains(out, "condition was not satisfied") {
		t.Errorf("diagnostic should describe the failure, got: %s", out)
	}
	if !strings.Contains(out, "30ms") && !strings.Contains(out, "30m") {
		t.Errorf("diagnostic should mention the timeout window, got: %s", out)
	}
}

// TestEventually_ZeroTimeoutChecksOnce: timeout <= 0 with a true predicate
// passes after exactly one invocation.
func TestEventually_ZeroTimeoutChecksOnce(t *testing.T) {
	calls := 0
	Eventually(t, func() bool {
		calls++
		return true
	}, 0, 0)

	if calls != 1 {
		t.Errorf("zero timeout should check predicate exactly once, got %d calls", calls)
	}
	if t.Failed() {
		t.Error("should not fail when single check returns true")
	}
}

// TestEventually_ZeroTimeoutChecksOnce_Negative: negative timeout is also
// treated as "check once" per the timeout <= 0 branch.
func TestEventually_ZeroTimeoutChecksOnce_Negative(t *testing.T) {
	calls := 0
	Eventually(t, func() bool {
		calls++
		return true
	}, -1*time.Second, 0)

	if calls != 1 {
		t.Errorf("negative timeout should check predicate exactly once, got %d", calls)
	}
}

// TestEventually_ZeroTimeoutFailsImmediately: timeout <= 0 with a false
// predicate fails immediately with the single-check diagnostic.
func TestEventually_ZeroTimeoutFailsImmediately(t *testing.T) {
	calls := 0
	ftb := runEventuallyInGoroutine(
		func() bool {
			calls++
			return false
		},
		0,
		0,
	)

	if calls != 1 {
		t.Errorf("zero-timeout failure path should call predicate once, got %d", calls)
	}
	if !ftb.Failed() {
		t.Fatal("zero timeout with false predicate should fail")
	}
	out := string(ftb.Buffer())
	if !strings.Contains(out, "checked once") {
		t.Errorf("diagnostic should mention single-check semantics, got: %s", out)
	}
}

// TestEventually_NilPredicateFails: nil predicate fails fast with the
// "predicate must not be nil" diagnostic.
func TestEventually_NilPredicateFails(t *testing.T) {
	ftb := runEventuallyInGoroutine(nil, time.Second, time.Millisecond)

	if !ftb.Failed() {
		t.Fatal("nil predicate should fail")
	}
	out := string(ftb.Buffer())
	if !strings.Contains(out, "predicate must not be nil") {
		t.Errorf("diagnostic should explain nil predicate, got: %s", out)
	}
}

// TestEventually_DefaultIntervalWhenZero: interval <= 0 falls back to the
// default interval and the predicate still converges.
func TestEventually_DefaultIntervalWhenZero(t *testing.T) {
	calls := 0
	Eventually(t, func() bool {
		calls++
		return calls >= 3
	}, 200*time.Millisecond, 0)

	if calls < 3 {
		t.Errorf("should have converged with default interval, got %d calls", calls)
	}
	if t.Failed() {
		t.Error("should not fail when predicate succeeds with default interval")
	}
}

// TestEventually_NegativeIntervalUsesDefault: interval < 0 is also treated as
// "use default" by the interval <= 0 branch.
func TestEventually_NegativeIntervalUsesDefault(t *testing.T) {
	calls := 0
	Eventually(t, func() bool {
		calls++
		return true
	}, 100*time.Millisecond, -5*time.Millisecond)

	if calls < 1 {
		t.Errorf("should have run at least once with default interval, got %d", calls)
	}
	if t.Failed() {
		t.Error("should not fail when predicate succeeds with negative interval")
	}
}
