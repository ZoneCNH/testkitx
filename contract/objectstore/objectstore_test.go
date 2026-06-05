package objectstore

import (
	"context"
	"io"
	"strings"
	"testing"
)

func TestRunObjectStore(t *testing.T) {
	t.Parallel()
	RunObjectStore(t, fakeStore{})
}

type fakeStore struct{}

func (fakeStore) Put(context.Context, Object) error { return nil }
func (fakeStore) Get(context.Context, string, string) (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader("value")), nil
}
func (fakeStore) Delete(context.Context, string, string) error { return nil }
