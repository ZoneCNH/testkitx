package kv

import (
	"context"

	"github.com/ZoneCNH/testkitx/requirex"
)

type Store interface {
	Set(context.Context, string, []byte) error
	Get(context.Context, string) ([]byte, error)
	Delete(context.Context, string) error
}

func RunBasic(t requirex.TestingT, store Store) {
	t.Helper()
	ctx := context.Background()
	key := "testkitx.contract.kv.basic"
	value := []byte("value")
	requirex.NoError(t, store.Set(ctx, key, value))
	got, err := store.Get(ctx, key)
	requirex.NoError(t, err)
	requirex.Equal(t, string(value), string(got))
	requirex.NoError(t, store.Delete(ctx, key))
}
