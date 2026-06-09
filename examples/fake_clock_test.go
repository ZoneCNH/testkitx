// Package examples demonstrates how to use testkitx helpers for deterministic
// testing. This file specifically shows the FakeClock from clocktest.
//
// Deterministic clocks eliminate flaky tests caused by real wall-clock drift,
// CI runner pauses, or timezone differences. By controlling time explicitly
// tests become reproducible and fast.
package examples

import (
	"testing"
	"time"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/clocktest"
)

// TestFakeClockAdvance shows how to advance a fake clock to a known point
// in time and assert the result.
func TestFakeClockAdvance(t *testing.T) {
	t.Parallel()

	// Create a clock pinned to a fixed instant.
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	clock := clocktest.New(start)

	// Advance by exactly 5 seconds — no real wall time passes.
	clock.Advance(5 * time.Second)

	want := start.Add(5 * time.Second)
	if got := clock.Now(); !got.Equal(want) {
		t.Fatalf("after advance: Now() = %s, want %s", got, want)
	}
}

// TestFakeClockMultipleAdvances demonstrates chaining multiple advances.
func TestFakeClockMultipleAdvances(t *testing.T) {
	t.Parallel()

	clock := clocktest.New(time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC))

	// First advance: move to lunch.
	clock.Advance(1 * time.Hour)
	// Second advance: end of workday.
	clock.Advance(7 * time.Hour)

	want := time.Date(2026, 6, 1, 20, 0, 0, 0, time.UTC)
	if got := clock.Now(); !got.Equal(want) {
		t.Fatalf("after two advances: Now() = %s, want %s", got, want)
	}
}

// TestFakeClockWithTimer shows a pattern for testing time-dependent logic
// using the fake clock's Now() as a deadline proxy.
func TestFakeClockWithTimer(t *testing.T) {
	t.Parallel()

	clock := clocktest.New(time.Date(2026, 3, 15, 10, 0, 0, 0, time.UTC))
	deadline := clock.Now().Add(30 * time.Second)

	// Simulate time passing in steps.
	clock.Advance(10 * time.Second)
	if clock.Now().After(deadline) {
		t.Fatal("deadline should not have passed yet")
	}

	clock.Advance(25 * time.Second)
	if !clock.Now().After(deadline) {
		t.Fatal("deadline should have passed now")
	}
}
