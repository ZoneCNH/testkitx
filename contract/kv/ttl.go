package kv

import (
	"context"
	"time"

	"github.com/ZoneCNH/testkitx/requirex"
)

type TTLStore interface {
	Store
	SetTTL(context.Context, string, []byte, time.Duration) error
}

func RunTTL(t requirex.TestingT, store TTLStore) {
	t.Helper()
	requirex.NoError(t, store.SetTTL(context.Background(), "testkitx.contract.kv.ttl", []byte("value"), time.Minute))
}
