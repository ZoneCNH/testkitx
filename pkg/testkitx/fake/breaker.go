package fake

import (
	"sync"
)

// BreakerState represents the state of a circuit breaker.
type BreakerState int

const (
	BreakerClosed   BreakerState = 0
	BreakerOpen     BreakerState = 1
	BreakerHalfOpen BreakerState = 2
)

func (s BreakerState) String() string {
	switch s {
	case BreakerClosed:
		return "closed"
	case BreakerOpen:
		return "open"
	case BreakerHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// Breaker mirrors the resiliencx.Breaker interface consumed by FoundationX modules.
type Breaker interface {
	State() BreakerState
	Allow() bool
	RecordSuccess()
	RecordFailure()
}

// FakeBreakerImpl is a deterministic fake circuit breaker. It implements Breaker.
type FakeBreakerImpl struct {
	mu    sync.Mutex
	state BreakerState
}

// Compile-time contract: *FakeBreakerImpl implements Breaker.
var _ Breaker = (*FakeBreakerImpl)(nil)

// FakeBreaker creates a deterministic fake circuit breaker with the given initial state.
func FakeBreaker(initial BreakerState) Breaker {
	return &FakeBreakerImpl{state: initial}
}

// State returns the current breaker state.
func (b *FakeBreakerImpl) State() BreakerState {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}

// Allow returns true when the breaker is not open.
func (b *FakeBreakerImpl) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state != BreakerOpen
}

// RecordSuccess transitions the breaker to closed.
func (b *FakeBreakerImpl) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.state = BreakerClosed
}

// RecordFailure transitions the breaker to open.
func (b *FakeBreakerImpl) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.state = BreakerOpen
}

// SetState allows tests to directly set the breaker state.
func (b *FakeBreakerImpl) SetState(state BreakerState) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.state = state
}
