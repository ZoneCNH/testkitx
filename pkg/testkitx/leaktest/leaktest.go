// Package leaktest provides lightweight goroutine leak checks for focused tests.
package leaktest

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

type Snapshot struct{ Goroutines int }

func Capture() Snapshot { return Snapshot{Goroutines: runtime.NumGoroutine()} }
func Check(start Snapshot, tolerance int) error {
	for i := 0; i < 20; i++ {
		runtime.GC()
		time.Sleep(time.Millisecond)
	}
	current := runtime.NumGoroutine()
	if current > start.Goroutines+tolerance {
		return fmt.Errorf("goroutine leak: before=%d after=%d tolerance=%d", start.Goroutines, current, tolerance)
	}
	return nil
}
func RequireNoLeak(t testing.TB, start Snapshot, tolerance int) {
	t.Helper()
	if err := Check(start, tolerance); err != nil {
		t.Fatal(err)
	}
}
