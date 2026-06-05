package eventlog

import (
	"context"
	"testing"
)

func TestRunners(t *testing.T) {
	t.Parallel()
	log := fakeLog{}
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{name: "producer", run: func(t *testing.T) { RunProducer(t, log) }},
		{name: "consumer", run: func(t *testing.T) { RunConsumer(t, log) }},
		{name: "offset_commit", run: func(t *testing.T) { RunOffsetCommit(t, log) }},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) { t.Parallel(); tc.run(t) })
	}
}

type fakeLog struct{}

func (fakeLog) Append(_ context.Context, _ string, events []Event) ([]Event, error) {
	events[0].Offset = 1
	return events, nil
}
func (fakeLog) Read(context.Context, string, int64, int) ([]Event, error) {
	return []Event{{Offset: 1}}, nil
}
func (fakeLog) CommitOffset(context.Context, string, int64) error { return nil }
