package kv

import (
	"context"
	"testing"
	"time"
)

func TestRunners(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		run  func(*testing.T, fakeStore)
	}{
		{name: "basic", run: func(t *testing.T, store fakeStore) { RunBasic(t, store) }},
		{name: "ttl", run: func(t *testing.T, store fakeStore) { RunTTL(t, store) }},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.run(t, fakeStore{values: map[string][]byte{}})
		})
	}
}

type fakeStore struct{ values map[string][]byte }

func (s fakeStore) Set(_ context.Context, key string, value []byte) error {
	s.values[key] = append([]byte(nil), value...)
	return nil
}
func (s fakeStore) Get(_ context.Context, key string) ([]byte, error) {
	return append([]byte(nil), s.values[key]...), nil
}
func (s fakeStore) Delete(_ context.Context, key string) error { delete(s.values, key); return nil }
func (s fakeStore) SetTTL(ctx context.Context, key string, value []byte, _ time.Duration) error {
	return s.Set(ctx, key, value)
}
