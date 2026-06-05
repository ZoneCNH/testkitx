package timeseries

import (
	"context"

	"github.com/ZoneCNH/testkitx/requirex"
)

type StableQuerier interface {
	Stable(context.Context, string) (bool, error)
}

func RunStable(t requirex.TestingT, querier StableQuerier) {
	t.Helper()
	stable, err := querier.Stable(context.Background(), "testkitx.contract.timeseries")
	requirex.NoError(t, err)
	if !stable {
		t.Fatalf("expected timeseries stability signal")
	}
}
