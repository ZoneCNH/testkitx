package fake

import (
	"sync"
	"time"
)

// FakeClock is a deterministic clock that does not call time.Now() or
// math.Rand(). The current time only advances via Advance or Set.
type FakeClock struct {
	mu  sync.Mutex
	now time.Time
}

// FakeClock creates a deterministic clock whose Now() returns at.
func Clock(at time.Time) *FakeClock {
	return &FakeClock{now: at}
}

// Now returns the current fake time. It never calls time.Now().
func (c *FakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

// Advance moves the clock forward by d and returns the new time.
func (c *FakeClock) Advance(d time.Duration) time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
	return c.now
}

// Set moves the clock to t and returns the new time.
func (c *FakeClock) Set(t time.Time) time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = t
	return c.now
}
