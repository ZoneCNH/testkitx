package eventlog

import (
	"context"

	"github.com/ZoneCNH/testkitx/requirex"
)

type Consumer interface {
	Read(context.Context, string, int64, int) ([]Event, error)
}

func RunConsumer(t requirex.TestingT, consumer Consumer) {
	t.Helper()
	events, err := consumer.Read(context.Background(), "testkitx.contract.eventlog", 0, 1)
	requirex.NoError(t, err)
	if events == nil {
		t.Fatalf("expected non-nil event slice")
	}
}
