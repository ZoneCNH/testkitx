package requirex

import (
	"errors"
	"testing"
	"time"
)

type typedErr struct{ kind string }

func (e typedErr) Error() string     { return e.kind }
func (e typedErr) ErrorKind() string { return e.kind }

type kindMethodErr struct{ kind string }

func (e kindMethodErr) Error() string { return e.kind }
func (e kindMethodErr) Kind() string  { return e.kind }

type fieldKindErr struct{ Kind string }

func (e *fieldKindErr) Error() string { return e.Kind }

func TestNoErrorAcceptsNil(t *testing.T) {
	t.Parallel()
	NoError(t, nil)
}

func TestEqualAcceptsMatchingComparable(t *testing.T) {
	t.Parallel()
	Equal(t, "want", "want")
}

func TestErrorKindAcceptsTypedError(t *testing.T) {
	t.Parallel()
	ErrorKind(t, typedErr{kind: "validation"}, "validation")
}

func TestErrorKindOneOfAcceptsAnyWantedKind(t *testing.T) {
	t.Parallel()
	ErrorKindOneOf(t, typedErr{kind: "timeout"}, "validation", "timeout")
}

func TestErrorKindAcceptsSupportedKindShapes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		err  error
	}{
		{name: "error-kind-method", err: typedErr{kind: "validation"}},
		{name: "kind-method", err: kindMethodErr{kind: "timeout"}},
		{name: "kind-field", err: &fieldKindErr{Kind: "canceled"}},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ErrorKind(t, tc.err, tc.err.Error())
		})
	}
}

func TestNoSecretLeakIgnoresMaskedValues(t *testing.T) {
	t.Parallel()
	NoSecretLeak(t, struct{ Secret string }{Secret: "***"}, "plain-text")
}

func TestEventuallyWaitsUntilPredicatePasses(t *testing.T) {
	t.Parallel()
	attempts := 0
	Eventually(t, time.Second, time.Millisecond, func() bool {
		attempts++
		return attempts == 2
	})
}

func TestNoGoroutineLeakAcceptsNoGrowth(t *testing.T) {
	t.Parallel()
	NoGoroutineLeak(t, 10, 10)
}

func TestErrorKindRejectsUntypedErrorsViaHelperHarness(t *testing.T) {
	t.Parallel()
	h := &harness{}
	ErrorKind(h, errors.New("plain"), "validation")
	if !h.failed {
		t.Fatal("expected helper harness to record failure")
	}
}

type harness struct{ failed bool }

func (h *harness) Helper()               {}
func (h *harness) Fatalf(string, ...any) { h.failed = true }
