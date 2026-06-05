package sql

import (
	"context"

	"github.com/ZoneCNH/testkitx/requirex"
)

type Pool interface {
	Stats(context.Context) (PoolStats, error)
}

type PoolStats struct {
	Open  int
	InUse int
	Idle  int
}

func RunPool(t requirex.TestingT, pool Pool) {
	t.Helper()
	stats, err := pool.Stats(context.Background())
	requirex.NoError(t, err)
	if stats.Open < 0 || stats.InUse < 0 || stats.Idle < 0 {
		t.Fatalf("pool stats must be non-negative: %+v", stats)
	}
}
