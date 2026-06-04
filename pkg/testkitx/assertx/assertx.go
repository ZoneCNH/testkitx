// Package assertx provides small deterministic assertions for testkitx users.
package assertx

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func Equal[T any](t testing.TB, want, got T) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("values differ\nwant: %#v\ngot:  %#v", want, got)
	}
}

func NoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func ErrorIs(t testing.TB, err, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("expected error %v to match %v", err, target)
	}
}

func Eventually(t testing.TB, timeout, interval time.Duration, check func() error) {
	t.Helper()
	if timeout <= 0 || interval <= 0 {
		t.Fatalf("timeout and interval must be positive, got %s/%s", timeout, interval)
	}
	deadline := time.Now().Add(timeout)
	var last error
	for {
		if err := check(); err == nil {
			return
		} else {
			last = err
		}
		if time.Now().After(deadline) {
			t.Fatalf("condition not met within %s: %v", timeout, last)
		}
		time.Sleep(interval)
	}
}

func Failf(format string, args ...any) error { return fmt.Errorf(format, args...) }
