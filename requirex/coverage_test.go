package requirex

import (
	"errors"
	"testing"
	"time"
)

// mockTB implements TestingT without calling runtime.Goexit on Fatalf.
type mockTB struct {
	failed bool
}

func (m *mockTB) Helper()                           {}
func (m *mockTB) Fatalf(format string, args ...any) { m.failed = true }

func TestEqualMismatch(t *testing.T) {
	h := &mockTB{}
	Equal(h, "want", "got")
	if !h.failed {
		t.Fatal("expected failure on mismatch")
	}
}

func TestNoErrorNonNil(t *testing.T) {
	h := &mockTB{}
	NoError(h, errors.New("boom"))
	if !h.failed {
		t.Fatal("expected failure on non-nil error")
	}
}

func TestErrorKindNilError(t *testing.T) {
	h := &mockTB{}
	ErrorKind(h, nil, "validation")
	if !h.failed {
		t.Fatal("expected failure on nil error")
	}
}

func TestErrorKindMismatch(t *testing.T) {
	h := &mockTB{}
	ErrorKind(h, typedErr{kind: "timeout"}, "validation")
	if !h.failed {
		t.Fatal("expected failure on kind mismatch")
	}
}

func TestErrorKindOneOfNilError(t *testing.T) {
	h := &mockTB{}
	ErrorKindOneOf(h, nil, "a", "b")
	if !h.failed {
		t.Fatal("expected failure on nil error")
	}
}

func TestErrorKindOneOfMismatch(t *testing.T) {
	h := &mockTB{}
	ErrorKindOneOf(h, typedErr{kind: "timeout"}, "validation", "auth")
	if !h.failed {
		t.Fatal("expected failure on kind mismatch")
	}
}

func TestErrorKindOneOfEmptyWants(t *testing.T) {
	h := &mockTB{}
	ErrorKindOneOf(h, typedErr{kind: "x"})
	if !h.failed {
		t.Fatal("expected failure on empty wants")
	}
}

func TestEventuallyNilPredicate(t *testing.T) {
	h := &panicMockTB{}
	defer func() { _ = recover() }()
	Eventually(h, time.Second, time.Millisecond, nil)
	if !h.failed {
		t.Fatal("expected failure on nil predicate")
	}
}

// panicMockTB panics on Fatalf to stop execution (like real testing.T).
type panicMockTB struct {
	failed bool
}

func (m *panicMockTB) Helper()                           {}
func (m *panicMockTB) Fatalf(format string, args ...any) { m.failed = true; panic("fatal") }

func TestEventuallyZeroTimeout(t *testing.T) {
	h := &panicMockTB{}
	defer func() { _ = recover() }()
	Eventually(h, 0, time.Millisecond, func() bool { return true })
	if !h.failed {
		t.Fatal("expected failure on zero timeout")
	}
}

func TestEventuallyNegativeTimeout(t *testing.T) {
	h := &panicMockTB{}
	defer func() { _ = recover() }()
	Eventually(h, -time.Second, time.Millisecond, func() bool { return true })
	if !h.failed {
		t.Fatal("expected failure on negative timeout")
	}
}

func TestEventuallyTimeoutExceeded(t *testing.T) {
	h := &panicMockTB{}
	defer func() { _ = recover() }()
	Eventually(h, 10*time.Millisecond, time.Millisecond, func() bool { return false })
	if !h.failed {
		t.Fatal("expected failure on timeout exceeded")
	}
}

func TestEventuallyZeroIntervalDefaults(t *testing.T) {
	attempts := 0
	Eventually(t, 100*time.Millisecond, 0, func() bool {
		attempts++
		return attempts >= 2
	})
}

func TestNoGoroutineLeakDetected(t *testing.T) {
	h := &mockTB{}
	NoGoroutineLeak(h, 5, 10)
	if !h.failed {
		t.Fatal("expected failure on leak")
	}
}

func TestNoSecretLeakDetected(t *testing.T) {
	h := &mockTB{}
	NoSecretLeak(h, struct{ Key string }{Key: "my-secret"}, "my-secret")
	if !h.failed {
		t.Fatal("expected failure on secret leak")
	}
}

func TestNoSecretLeakIgnoresEmptySecret(t *testing.T) {
	NoSecretLeak(t, "anything", "")
}

func TestErrorKindFieldReflection(t *testing.T) {
	ErrorKind(t, &fieldKindErr{Kind: "canceled"}, "canceled")
}
