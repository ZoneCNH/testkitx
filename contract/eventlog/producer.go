package eventlog

import (
	"context"

	"github.com/ZoneCNH/testkitx/requirex"
)

type Event struct {
	Stream   string
	Offset   int64
	Payload  []byte
	Metadata map[string]string
}

type Producer interface {
	Append(context.Context, string, []Event) ([]Event, error)
}

func RunProducer(t requirex.TestingT, producer Producer) {
	t.Helper()
	events, err := producer.Append(context.Background(), "testkitx.contract.eventlog", []Event{{Payload: []byte("value")}})
	requirex.NoError(t, err)
	if len(events) == 0 {
		t.Fatalf("expected appended events")
	}
}
