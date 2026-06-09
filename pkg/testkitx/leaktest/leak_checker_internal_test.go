package leaktest

import (
	"testing"
)

// probeT is a testing.TB that captures failures without calling runtime.Goexit.
type probeT struct {
	testing.TB
	failed bool
}

func (p *probeT) Helper()                           {}
func (p *probeT) Errorf(format string, args ...any) { p.failed = true }
func (p *probeT) Fatalf(format string, args ...any) { p.failed = true }
func (p *probeT) Error(args ...any)                 { p.failed = true }
func (p *probeT) Fatal(args ...any)                 { p.failed = true }
func (p *probeT) FailNow()                          { p.failed = true }
func (p *probeT) Failed() bool                      { return p.failed }

func TestCheckLeakCleanupDetectsLeak(t *testing.T) {
	// Directly test the extracted cleanup function with a probeT.
	// After will be > 0, before is 0, so cleanup detects a "leak" and calls Errorf.
	before := Snapshot{Goroutines: 0}
	p := &probeT{TB: t}
	checkLeakCleanup(p, before, nil)
	if !p.failed {
		t.Fatal("expected cleanup to detect leak")
	}
}

func TestCheckLeakCleanupWithIgnorePatterns(t *testing.T) {
	// Use an empty stack scenario: set before high enough that after <= before.
	// This tests the ignore patterns path without triggering leak detection.
	// Actually, we want to test the filtering. Use a real scenario:
	// before=0, after>0, but patterns match everything.
	// The goroutine dump includes blocks like "goroutine N [state]:\n ... stack ..."
	// Use a very broad pattern that matches all blocks.
	before := Snapshot{Goroutines: 0}
	p := &probeT{TB: t}
	// Use a pattern that matches the runtime stack trace format.
	checkLeakCleanup(p, before, []string{"goroutine"})
	if p.failed {
		t.Fatal("expected no leak with 'goroutine' ignore pattern matching all blocks")
	}
}

func TestCheckLeakCleanupNoLeak(t *testing.T) {
	// Use a very high goroutine count so after <= before → early return.
	before := Snapshot{Goroutines: 99999}
	p := &probeT{TB: t}
	checkLeakCleanup(p, before, nil)
	if p.failed {
		t.Fatal("expected no leak with high before count")
	}
}
