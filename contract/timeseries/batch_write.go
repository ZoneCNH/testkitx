package timeseries

import (
	"context"
	"time"

	"github.com/ZoneCNH/testkitx/requirex"
)

type BatchWriter interface {
	WriteBatch(context.Context, []Point) error
}

func RunBatchWrite(t requirex.TestingT, writer BatchWriter) {
	t.Helper()
	requirex.NoError(t, writer.WriteBatch(context.Background(), []Point{{Metric: "testkitx.contract.timeseries", Time: time.Now().UTC(), Value: 1}}))
}
