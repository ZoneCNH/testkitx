package timeseries

import (
	"context"
	"time"

	"github.com/ZoneCNH/testkitx/requirex"
)

type ChildTableCreator interface {
	CreateChildTable(context.Context, string, time.Time) error
}

func RunChildTable(t requirex.TestingT, creator ChildTableCreator) {
	t.Helper()
	requirex.NoError(t, creator.CreateChildTable(context.Background(), "testkitx_contract_timeseries", time.Now().UTC()))
}
