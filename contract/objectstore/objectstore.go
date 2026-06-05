package objectstore

import (
	"context"
	"io"
	"strings"

	"github.com/ZoneCNH/testkitx/requirex"
)

type Object struct {
	Bucket   string
	Key      string
	Body     io.Reader
	Metadata map[string]string
}
type Store interface {
	Put(context.Context, Object) error
	Get(context.Context, string, string) (io.ReadCloser, error)
	Delete(context.Context, string, string) error
}

func RunObjectStore(t requirex.TestingT, store Store) {
	t.Helper()
	ctx := context.Background()
	bucket, key := "testkitx-contract", "object"
	requirex.NoError(t, store.Put(ctx, Object{Bucket: bucket, Key: key, Body: strings.NewReader("value")}))
	body, err := store.Get(ctx, bucket, key)
	requirex.NoError(t, err)
	if body == nil {
		t.Fatalf("expected object body")
	}
	requirex.NoError(t, body.Close())
	requirex.NoError(t, store.Delete(ctx, bucket, key))
}
