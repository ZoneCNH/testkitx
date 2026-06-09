package timeseries

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRunners(t *testing.T) {
	t.Parallel()
	store := fakeStore{}
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{name: "timeseries", run: func(t *testing.T) { RunTimeSeries(t, store) }},
		{name: "stable", run: func(t *testing.T) { RunStable(t, store) }},
		{name: "child_table", run: func(t *testing.T) { RunChildTable(t, store) }},
		{name: "batch_write", run: func(t *testing.T) { RunBatchWrite(t, store) }},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) { t.Parallel(); tc.run(t) })
	}
}

type fakeStore struct{}

func (fakeStore) Write(context.Context, Point) error { return nil }
func (fakeStore) Query(context.Context, string, time.Time, time.Time) ([]Point, error) {
	return []Point{}, nil
}
func (fakeStore) Stable(context.Context, string) (bool, error)              { return true, nil }
func (fakeStore) CreateChildTable(context.Context, string, time.Time) error { return nil }
func (fakeStore) WriteBatch(context.Context, []Point) error                 { return nil }

func TestRunStableRequiresStableSignal(t *testing.T) {
	t.Parallel()
	probe := &probeT{}
	RunStable(probe, unstableStore{})
	if !probe.failed {
		t.Fatalf("expected RunStable to fail when Stable returns false")
	}
}

func TestRunStableUsesContractMetric(t *testing.T) {
	t.Parallel()
	store := &recordingStableStore{}
	RunStable(t, store)
	if store.metric != "testkitx.contract.timeseries" {
		t.Fatalf("expected contract metric, got %q", store.metric)
	}
}

func TestRunStableFailsOnStableError(t *testing.T) {
	t.Parallel()
	probe := &probeT{}
	RunStable(probe, errorStableStore{})
	if !probe.failed {
		t.Fatalf("expected RunStable to fail when Stable returns an error")
	}
}

type unstableStore struct{ fakeStore }

func (unstableStore) Stable(context.Context, string) (bool, error) { return false, nil }

type recordingStableStore struct {
	fakeStore
	metric string
}

func (s *recordingStableStore) Stable(_ context.Context, metric string) (bool, error) {
	s.metric = metric
	return true, nil
}

type errorStableStore struct{ fakeStore }

var errStable = errors.New("stable failed")

func (errorStableStore) Stable(context.Context, string) (bool, error) { return false, errStable }

type probeT struct{ failed bool }

func (p *probeT) Helper()               {}
func (p *probeT) Fatalf(string, ...any) { p.failed = true }

func TestRunTimeSeriesNilPoints(t *testing.T) {
	t.Parallel()
	probe := &probeT{}
	RunTimeSeries(probe, nilPointsStore{})
	if !probe.failed {
		t.Fatal("expected failure for nil points")
	}
}

type nilPointsStore struct{ fakeStore }

func (nilPointsStore) Query(context.Context, string, time.Time, time.Time) ([]Point, error) {
	return nil, nil
}
