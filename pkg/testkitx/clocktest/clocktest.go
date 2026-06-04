// Package clocktest provides a deterministic fake clock.
package clocktest

import (
	"sync"
	"time"
)

type Clock struct {
	mu  sync.Mutex
	now time.Time
}

func New(start time.Time) *Clock { return &Clock{now: start} }
func (c *Clock) Now() time.Time  { c.mu.Lock(); defer c.mu.Unlock(); return c.now }
func (c *Clock) Advance(d time.Duration) time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
	return c.now
}
