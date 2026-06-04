package clocktest_test

import (
	"testing"
	"time"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/clocktest"
)

func TestFakeClockAdvancesDeterministically(t *testing.T) {
	t.Parallel()
	start := time.Date(2026, 6, 4, 0, 0, 0, 0, time.UTC)
	clock := clocktest.New(start)
	clock.Advance(2 * time.Second)
	if got := clock.Now(); !got.Equal(start.Add(2 * time.Second)) {
		t.Fatalf("Now() = %s, want %s", got, start.Add(2*time.Second))
	}
}
