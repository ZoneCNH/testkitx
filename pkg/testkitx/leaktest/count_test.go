package leaktest

import (
	"testing"
)

func TestCountUnexpectedGoroutinesFilters(t *testing.T) {
	t.Parallel()
	stack := "goroutine 1 [running]:\ntest\n\n\nruntime.main()\n\ngoroutine 2 [sleep]:\nother\n"
	// Two non-empty blocks: "goroutine 1..." and "runtime.main()" is filtered by "runtime."
	// Actually blocks split by "\n\n": "goroutine 1 [running]:\ntest", "\n", "runtime.main()", "\n", "goroutine 2 [sleep]:\nother", "\n"
	// Non-empty trimmed blocks without "runtime.": "goroutine 1 [running]:\ntest", "goroutine 2 [sleep]:\nother"
	got := countUnexpectedGoroutines(stack, []string{"runtime."})
	if got != 2 {
		t.Fatalf("expected 2 unexpected, got %d", got)
	}
}

func TestCountUnexpectedGoroutinesAllFiltered(t *testing.T) {
	t.Parallel()
	stack := "runtime.main()\n\nruntime.gc()\n"
	got := countUnexpectedGoroutines(stack, []string{"runtime."})
	if got != 0 {
		t.Fatalf("expected 0 unexpected, got %d", got)
	}
}

func TestCountUnexpectedGoroutinesEmpty(t *testing.T) {
	t.Parallel()
	got := countUnexpectedGoroutines("", nil)
	if got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}
