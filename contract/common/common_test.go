package common

import (
	"context"
	"testing"
)

func TestTestIDsMatchExecutionPlan(t *testing.T) {
	t.Parallel()
	want := []string{
		LifecycleStartID,
		LifecyclePingID,
		LifecycleCloseID,
		ConfigInvalidID,
		ErrorStandardKindID,
		SecretNoLeakID,
		ResilienceCancelID,
		ObservabilityMetricsID,
	}
	if len(TestIDs) != len(want) {
		t.Fatalf("expected %d IDs, got %d", len(want), len(TestIDs))
	}
	for i, id := range want {
		if TestIDs[i] != id {
			t.Fatalf("id[%d]: expected %q, got %q", i, id, TestIDs[i])
		}
	}
}

func TestCommonRunners(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{name: LifecycleStartID, run: func(t *testing.T) { RunLifecycleStart(t, fakeFactory{}) }},
		{name: LifecyclePingID, run: func(t *testing.T) { RunLifecyclePing(t, fakeFactory{}) }},
		{name: LifecycleCloseID, run: func(t *testing.T) { RunLifecycleCloseIdempotent(t, fakeFactory{}) }},
		{name: ConfigInvalidID, run: func(t *testing.T) { RunInvalidConfig(t, fakeFactory{}) }},
		{name: ErrorStandardKindID, run: func(t *testing.T) {
			RunStandardErrorKind(t, Error{Kind: "validation", Op: "op", Message: "bad"}, "validation")
		}},
		{name: SecretNoLeakID, run: func(t *testing.T) { RunSecretNoLeak(t, Config{Name: "ok", Secret: "***"}, "plain-text") }},
		{name: ResilienceCancelID, run: func(t *testing.T) {
			RunContextCancel(t, func(ctx context.Context) error { return Error{Kind: "canceled", Message: ctx.Err().Error()} })
		}},
		{name: ObservabilityMetricsID, run: func(t *testing.T) { RunObservabilityMetrics(t, fakeMetrics{"requests": 1}, "requests") }},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.run(t)
		})
	}
}

type fakeFactory struct{}

func (fakeFactory) New(ctx context.Context, cfg Config) (Runtime, error) {
	if cfg.Name == "" {
		return nil, Error{Kind: "validation", Op: "new", Message: "name is required"}
	}
	return &fakeRuntime{}, ctx.Err()
}

type fakeRuntime struct{ closed bool }

func (r *fakeRuntime) Start(ctx context.Context) error { return ctx.Err() }
func (r *fakeRuntime) Ping(ctx context.Context) error  { return ctx.Err() }
func (r *fakeRuntime) Close(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.closed = true
	return nil
}

type fakeMetrics map[string]float64

func (m fakeMetrics) Metric(name string) (float64, bool) {
	value, ok := m[name]
	return value, ok
}
