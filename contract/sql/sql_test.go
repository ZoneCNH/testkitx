package sql

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"testing"
)

func TestRunners(t *testing.T) {
	t.Parallel()
	db := fakeDB{}
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{name: "exec_query", run: func(t *testing.T) { RunExecQuery(t, db) }},
		{name: "transaction", run: func(t *testing.T) { RunTransaction(t, db) }},
		{name: "pool", run: func(t *testing.T) { RunPool(t, db) }},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) { t.Parallel(); tc.run(t) })
	}
}

type fakeDB struct{}
type fakeResult struct{}
type fakeRows struct{}
type fakeTx struct{}

func (fakeDB) Exec(context.Context, string, ...any) (Result, error) { return fakeResult{}, nil }
func (fakeDB) Query(context.Context, string, ...any) (Rows, error)  { return fakeRows{}, nil }
func (fakeDB) BeginTx(context.Context) (Tx, error)                  { return fakeTx{}, nil }
func (fakeDB) Stats(context.Context) (PoolStats, error)             { return PoolStats{Open: 1, Idle: 1}, nil }
func (fakeResult) RowsAffected() int64                              { return 1 }
func (fakeRows) Next() bool                                         { return false }
func (fakeRows) Close() error                                       { return nil }
func (fakeTx) Exec(context.Context, string, ...any) (Result, error) { return fakeResult{}, nil }
func (fakeTx) Commit(context.Context) error                         { return nil }
func (fakeTx) Rollback(context.Context) error                       { return nil }

func TestRunTransactionCoversCommitAndRollback(t *testing.T) {
	t.Parallel()
	db := &recordingDB{}
	RunTransaction(t, db)
	if db.beginCount != 2 {
		t.Fatalf("expected two transactions, got %d", db.beginCount)
	}
	if db.commitCount != 1 {
		t.Fatalf("expected one commit, got %d", db.commitCount)
	}
	if db.rollbackCount != 1 {
		t.Fatalf("expected one rollback, got %d", db.rollbackCount)
	}
	wantOperations := []string{
		"begin",
		"exec:testkitx contract tx",
		"commit",
		"begin",
		"exec:testkitx contract tx rollback",
		"rollback",
	}
	if len(db.operations) != len(wantOperations) {
		t.Fatalf("expected operations %v, got %v", wantOperations, db.operations)
	}
	for i, want := range wantOperations {
		if db.operations[i] != want {
			t.Fatalf("operation[%d]: expected %q, got %q", i, want, db.operations[i])
		}
	}
}

func TestRunTransactionFailsOnExecError(t *testing.T) {
	t.Parallel()
	probe := &probeT{}
	RunTransaction(probe, failingTxDB{})
	if !probe.failed {
		t.Fatalf("expected RunTransaction to fail on tx exec error")
	}
}

func TestRunExecQueryNilRows(t *testing.T) {
	t.Parallel()
	probe := &fatalProbeT{}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() { _ = recover() }()
		RunExecQuery(probe, nilRowsDB{})
	}()
	wg.Wait()
	if !probe.failed {
		t.Fatal("expected failure for nil rows")
	}
}

func TestRunPoolNegativeStats(t *testing.T) {
	t.Parallel()
	probe := &fatalProbeT{}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() { _ = recover() }()
		RunPool(probe, negStatsPool{})
	}()
	wg.Wait()
	if !probe.failed {
		t.Fatal("expected failure for negative stats")
	}
}

type recordingDB struct {
	beginCount    int
	commitCount   int
	rollbackCount int
	operations    []string
}

func (db *recordingDB) BeginTx(context.Context) (Tx, error) {
	db.beginCount++
	db.operations = append(db.operations, "begin")
	return &recordingTx{db: db}, nil
}

type recordingTx struct{ db *recordingDB }

func (tx *recordingTx) Exec(_ context.Context, query string, _ ...any) (Result, error) {
	tx.db.operations = append(tx.db.operations, "exec:"+query)
	return fakeResult{}, nil
}
func (tx *recordingTx) Commit(context.Context) error {
	tx.db.commitCount++
	tx.db.operations = append(tx.db.operations, "commit")
	return nil
}
func (tx *recordingTx) Rollback(context.Context) error {
	tx.db.rollbackCount++
	tx.db.operations = append(tx.db.operations, "rollback")
	return nil
}

type failingTxDB struct{}

func (failingTxDB) BeginTx(context.Context) (Tx, error) { return failingTx{}, nil }

type failingTx struct{}

var errTxExec = errors.New("tx exec failed")

func (failingTx) Exec(context.Context, string, ...any) (Result, error) {
	return fakeResult{}, errTxExec
}
func (failingTx) Commit(context.Context) error   { return nil }
func (failingTx) Rollback(context.Context) error { return nil }

// probeT is a simple test probe that records Fatalf calls without stopping execution.
type probeT struct{ failed bool }

func (p *probeT) Helper()               {}
func (p *probeT) Fatalf(string, ...any) { p.failed = true }

// fatalProbeT calls runtime.Goexit() on Fatalf, like *testing.T.
// Must be used in a goroutine with recover() to avoid test failures.
type fatalProbeT struct{ failed bool }

func (p *fatalProbeT) Helper()               {}
func (p *fatalProbeT) Fatalf(string, ...any) { p.failed = true; runtime.Goexit() }

type nilRowsDB struct{ fakeDB }

func (nilRowsDB) Query(context.Context, string, ...any) (Rows, error) { return nil, nil }

type negStatsPool struct{ fakeDB }

func (negStatsPool) Stats(context.Context) (PoolStats, error) {
	return PoolStats{Open: -1}, nil
}

func TestRunTransactionNilTx(t *testing.T) {
	t.Parallel()
	probe := &probeT{}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() { _ = recover() }()
		RunTransaction(probe, nilTxDB{})
	}()
	wg.Wait()
	if !probe.failed {
		t.Fatal("expected failure for nil tx")
	}
}

type nilTxDB struct{}

func (nilTxDB) BeginTx(context.Context) (Tx, error) { return nil, nil }
