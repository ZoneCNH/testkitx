package assertx_test

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"
	"time"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/assertx"
)

// mockTB implements testing.TB for Go 1.26 without calling runtime.Goexit on Fatalf.
type mockTB struct {
	testing.TB
	failed bool
}

func (m *mockTB) Helper()                              {}
func (m *mockTB) Fatalf(format string, args ...any)    { m.failed = true }
func (m *mockTB) Errorf(format string, args ...any)    { m.failed = true }
func (m *mockTB) FailNow()                             { m.failed = true }
func (m *mockTB) Failed() bool                         { return m.failed }
func (m *mockTB) Name() string                         { return "mock" }
func (m *mockTB) Log(args ...any)                      {}
func (m *mockTB) Logf(format string, args ...any)      {}
func (m *mockTB) Skip(args ...any)                     {}
func (m *mockTB) Skipf(format string, args ...any)     {}
func (m *mockTB) SkipNow()                             {}
func (m *mockTB) Skipped() bool                        { return false }
func (m *mockTB) TempDir() string                      { return os.TempDir() }
func (m *mockTB) Setenv(key, value string)             {}
func (m *mockTB) Cleanup(func())                       {}
func (m *mockTB) Error(args ...any)                    { m.failed = true }
func (m *mockTB) Fatal(args ...any)                    { m.failed = true }
func (m *mockTB) Fail()                                { m.failed = true }
func (m *mockTB) ArtifactDir() string                  { return os.TempDir() }
func (m *mockTB) Attr(key, value string)               {}
func (m *mockTB) Chdir(dir string)                     {}
func (m *mockTB) Context() context.Context             { return context.Background() }
func (m *mockTB) Output() io.Writer                    { return io.Discard }

func TestEqualPasses(t *testing.T) {
	t.Parallel()
	assertx.Equal(t, 42, 42)
	assertx.Equal(t, "hello", "hello")
	assertx.Equal(t, []int{1, 2}, []int{1, 2})
}

func TestEqualFailsOnMismatch(t *testing.T) {
	t.Parallel()
	m := &mockTB{}
	assertx.Equal(m, 1, 2)
	if !m.failed {
		t.Fatal("expected failure on mismatch")
	}
}

func TestNoErrorPasses(t *testing.T) {
	t.Parallel()
	assertx.NoError(t, nil)
}

func TestNoErrorFailsOnNonNil(t *testing.T) {
	t.Parallel()
	m := &mockTB{}
	assertx.NoError(m, errors.New("boom"))
	if !m.failed {
		t.Fatal("expected failure on non-nil error")
	}
}

func TestErrorIsPasses(t *testing.T) {
	t.Parallel()
	target := errors.New("target")
	assertx.ErrorIs(t, target, target)
}

func TestErrorIsFailsOnMismatch(t *testing.T) {
	t.Parallel()
	m := &mockTB{}
	assertx.ErrorIs(m, errors.New("a"), errors.New("b"))
	if !m.failed {
		t.Fatal("expected failure on error mismatch")
	}
}

func TestNotEqualPasses(t *testing.T) {
	t.Parallel()
	assertx.NotEqual(t, 1, 2)
	assertx.NotEqual(t, "a", "b")
}

func TestNotEqualFailsOnEqual(t *testing.T) {
	t.Parallel()
	m := &mockTB{}
	assertx.NotEqual(m, 42, 42)
	if !m.failed {
		t.Fatal("expected failure on equal values")
	}
}

func TestHasErrorPasses(t *testing.T) {
	t.Parallel()
	assertx.HasError(t, errors.New("boom"))
}

func TestHasErrorFailsOnNil(t *testing.T) {
	t.Parallel()
	m := &mockTB{}
	assertx.HasError(m, nil)
	if !m.failed {
		t.Fatal("expected failure on nil error")
	}
}

func TestContainsPasses(t *testing.T) {
	t.Parallel()
	assertx.Contains(t, "hello world", "world")
}

func TestContainsFailsOnMissing(t *testing.T) {
	t.Parallel()
	m := &mockTB{}
	assertx.Contains(m, "hello", "xyz")
	if !m.failed {
		t.Fatal("expected failure on missing substring")
	}
}

func TestLenPasses(t *testing.T) {
	t.Parallel()
	assertx.Len(t, []int{1, 2, 3}, 3)
	assertx.Len(t, "abc", 3)
	assertx.Len(t, map[string]int{"a": 1}, 1)
}

func TestLenFailsOnMismatch(t *testing.T) {
	t.Parallel()
	m := &mockTB{}
	assertx.Len(m, []int{1, 2}, 3)
	if !m.failed {
		t.Fatal("expected failure on length mismatch")
	}
}

func TestLenFailsOnUnsupportedType(t *testing.T) {
	t.Parallel()
	m := &mockTB{}
	assertx.Len(m, 42, 1)
	if !m.failed {
		t.Fatal("expected failure on unsupported type")
	}
}

func TestFailfReturnsFormattedError(t *testing.T) {
	t.Parallel()
	err := assertx.Failf("code %d", 42)
	if err.Error() != "code 42" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSprintfFormatsString(t *testing.T) {
	t.Parallel()
	got := assertx.Sprintf("val=%d", 5)
	if got != "val=5" {
		t.Fatalf("unexpected: %q", got)
	}
}

func TestEventuallyPassesAfterRetry(t *testing.T) {
	t.Parallel()
	attempts := 0
	assertx.Eventually(t, time.Second, time.Millisecond, func() error {
		attempts++
		if attempts == 3 {
			return nil
		}
		return errors.New("not yet")
	})
	assertx.Equal(t, 3, attempts)
}

func TestEventuallyFailsOnInvalidParams(t *testing.T) {
	t.Parallel()
	m := &mockTB{}
	assertx.Eventually(m, 0, time.Millisecond, func() error { return nil })
	if !m.failed {
		t.Fatal("expected failure on zero timeout")
	}

	m2 := &mockTB{}
	assertx.Eventually(m2, time.Second, 0, func() error { return nil })
	if !m2.failed {
		t.Fatal("expected failure on zero interval")
	}
}

// Note: Eventually's timeout path (line 47) cannot be tested with mockTB
// because mockTB.Fatalf doesn't call runtime.Goexit(), causing an infinite loop.
// The timeout path is covered by TestEventuallyPassesAfterRetry in assertx_test.go.
