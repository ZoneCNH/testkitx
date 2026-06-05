package timeseries

import (
	"context"
	"time"

	"github.com/ZoneCNH/testkitx/requirex"
)

type Point struct {
	Metric string
	Time   time.Time
	Value  float64
	Tags   map[string]string
}
type Writer interface {
	Write(context.Context, Point) error
	Query(context.Context, string, time.Time, time.Time) ([]Point, error)
}

func RunTimeSeries(t requirex.TestingT, writer Writer) {
	t.Helper()
	now := time.Now().UTC()
	requirex.NoError(t, writer.Write(context.Background(), Point{Metric: "testkitx.contract.timeseries", Time: now, Value: 1}))
	points, err := writer.Query(context.Background(), "testkitx.contract.timeseries", now.Add(-time.Minute), now.Add(time.Minute))
	requirex.NoError(t, err)
	if points == nil {
		t.Fatalf("expected non-nil point slice")
	}
}
