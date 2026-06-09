package objectstore

import (
	"context"
	"io"
	"runtime"
	"strings"
	"sync"
	"testing"
)

func TestRunObjectStore(t *testing.T) {
	t.Parallel()
	RunObjectStore(t, fakeStore{})
}

func TestRunObjectStoreNilBody(t *testing.T) {
	t.Parallel()
	probe := &fatalProbeT{}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() { recover() }()
		RunObjectStore(probe, nilBodyStore{})
	}()
	wg.Wait()
	if !probe.failed {
		t.Fatal("expected failure for nil body")
	}
}

type fakeStore struct{}

func (fakeStore) Put(context.Context, Object) error { return nil }
func (fakeStore) Get(context.Context, string, string) (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader("value")), nil
}
func (fakeStore) Delete(context.Context, string, string) error { return nil }

type nilBodyStore struct{ fakeStore }

func (nilBodyStore) Get(context.Context, string, string) (io.ReadCloser, error) { return nil, nil }

type fatalProbeT struct{ failed bool }

func (p *fatalProbeT) Helper()               {}
func (p *fatalProbeT) Fatalf(string, ...any) { p.failed = true; runtime.Goexit() }
