package servicex

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestComposeStartsAndStopsInOrder(t *testing.T) {
	t.Parallel()
	var calls []string
	compose := NewCompose(fakeService{name: "a", calls: &calls}, fakeService{name: "b", calls: &calls})
	if err := compose.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	if err := compose.Stop(context.Background()); err != nil {
		t.Fatalf("stop: %v", err)
	}
	want := []string{"start:a", "start:b", "stop:b", "stop:a"}
	for i, call := range want {
		if calls[i] != call {
			t.Fatalf("call[%d]: expected %q, got %q", i, call, calls[i])
		}
	}
}

func TestWaitUntilReturnsWhenReady(t *testing.T) {
	t.Parallel()
	attempts := 0
	err := WaitUntil(context.Background(), time.Millisecond, func(context.Context) (bool, error) { attempts++; return attempts == 2, nil })
	if err != nil {
		t.Fatalf("wait: %v", err)
	}
}

func TestWaitUntilAddsDefaultDeadlineWhenCallerOmitsOne(t *testing.T) {
	t.Parallel()
	var hasDeadline bool
	err := WaitUntil(context.Background(), time.Millisecond, func(ctx context.Context) (bool, error) {
		_, hasDeadline = ctx.Deadline()
		return true, nil
	})
	if err != nil {
		t.Fatalf("wait: %v", err)
	}
	if !hasDeadline {
		t.Fatalf("expected WaitUntil to add a default deadline")
	}
}

func TestWaitUntilTreatsNilContextAsBackground(t *testing.T) {
	t.Parallel()
	var hasDeadline bool
	err := WaitUntil(nil, time.Millisecond, func(ctx context.Context) (bool, error) { //nolint:staticcheck // verifies the defensive nil-context branch.
		_, hasDeadline = ctx.Deadline()
		return true, nil
	})
	if err != nil {
		t.Fatalf("wait: %v", err)
	}
	if !hasDeadline {
		t.Fatalf("expected nil context to receive a default deadline")
	}
}

func TestWaitUntilHonorsCallerDeadline(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	err := WaitUntil(ctx, time.Millisecond, func(context.Context) (bool, error) { return false, nil })
	if !errors.Is(err, ErrWaitTimeout) || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected joined wait timeout and deadline exceeded, got %v", err)
	}
}

func TestWaitUntilReturnsCanceledContextBeforeReady(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	called := false
	err := WaitUntil(ctx, time.Millisecond, func(context.Context) (bool, error) {
		called = true
		return true, nil
	})
	if !errors.Is(err, ErrWaitTimeout) || !errors.Is(err, context.Canceled) {
		t.Fatalf("expected joined wait timeout and context canceled, got %v", err)
	}
	if called {
		t.Fatalf("expected canceled context to return before invoking readiness check")
	}
}

func TestWaitUntilRejectsNilReady(t *testing.T) {
	t.Parallel()
	if err := WaitUntil(context.Background(), time.Millisecond, nil); !errors.Is(err, ErrNilReady) {
		t.Fatalf("expected nil ready error, got %v", err)
	}
}

func TestWaitUntilPropagatesReadyError(t *testing.T) {
	t.Parallel()
	want := errors.New("ready failed")
	if err := WaitUntil(context.Background(), 0, func(context.Context) (bool, error) {
		return false, want
	}); !errors.Is(err, want) {
		t.Fatalf("expected ready error, got %v", err)
	}
}

func TestCheckHealthDelegates(t *testing.T) {
	t.Parallel()
	if err := CheckHealth(context.Background(), fakeHealth{}); err != nil {
		t.Fatalf("health: %v", err)
	}
}

type fakeService struct {
	name  string
	calls *[]string
}

func (s fakeService) Start(context.Context) error {
	*s.calls = append(*s.calls, "start:"+s.name)
	return nil
}
func (s fakeService) Stop(context.Context) error {
	*s.calls = append(*s.calls, "stop:"+s.name)
	return nil
}

type fakeHealth struct{}

func (fakeHealth) Healthy(context.Context) error { return nil }

func TestComposeStartFailsOnServiceError(t *testing.T) {
	t.Parallel()
	compose := NewCompose(failingService{}, failingService{startErr: true})
	err := compose.Start(context.Background())
	if err == nil || err.Error() != "start failed" {
		t.Fatalf("expected start error, got %v", err)
	}
}

func TestComposeStopFailsOnServiceError(t *testing.T) {
	t.Parallel()
	compose := NewCompose(failingService{}, failingService{stopErr: true})
	if err := compose.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	err := compose.Stop(context.Background())
	if err == nil || err.Error() != "stop failed" {
		t.Fatalf("expected stop error, got %v", err)
	}
}

func TestCheckHealthNilChecker(t *testing.T) {
	t.Parallel()
	err := CheckHealth(context.Background(), nil)
	if err == nil || err.Error() != "servicex: health checker is required" {
		t.Fatalf("expected nil checker error, got %v", err)
	}
}

type failingService struct {
	startErr bool
	stopErr  bool
}

func (s failingService) Start(context.Context) error {
	if s.startErr {
		return errors.New("start failed")
	}
	return nil
}

func (s failingService) Stop(context.Context) error {
	if s.stopErr {
		return errors.New("stop failed")
	}
	return nil
}
