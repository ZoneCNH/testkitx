package columnstore

import (
	"context"
	"testing"
)

func TestRunners(t *testing.T) {
	t.Parallel()
	store := fakeStore{}
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{name: "columnstore", run: func(t *testing.T) { RunColumnStore(t, store) }},
		{name: "batch_insert", run: func(t *testing.T) { RunBatchInsert(t, store) }},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) { t.Parallel(); tc.run(t) })
	}
}

func TestRunBatchInsertRowCheck(t *testing.T) {
	t.Parallel()
	probe := &probeT{}
	RunBatchInsert(probe, lowRowCountStore{})
	if !probe.failed {
		t.Fatal("expected failure for < 2 rows")
	}
}

func TestRunColumnStoreNilRows(t *testing.T) {
	t.Parallel()
	probe := &probeT{}
	RunColumnStore(probe, nilRowsStore{})
	if !probe.failed {
		t.Fatal("expected failure for nil rows")
	}
}

type fakeStore struct{}

func (fakeStore) Insert(context.Context, string, Row) error    { return nil }
func (fakeStore) Query(context.Context, string) ([]Row, error) { return []Row{{"id": "1"}}, nil }
func (fakeStore) InsertBatch(_ context.Context, _ string, rows []Row) (Result, error) {
	return Result{Rows: len(rows)}, nil
}

type lowRowCountStore struct{ fakeStore }

func (lowRowCountStore) InsertBatch(_ context.Context, _ string, _ []Row) (Result, error) {
	return Result{Rows: 1}, nil
}

type nilRowsStore struct{ fakeStore }

func (nilRowsStore) Query(context.Context, string) ([]Row, error) { return nil, nil }

type probeT struct{ failed bool }

func (p *probeT) Helper()               {}
func (p *probeT) Fatalf(string, ...any) { p.failed = true }
