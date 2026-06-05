package pubsub

import (
	"context"
	"testing"
)

func TestRunners(t *testing.T) {
	t.Parallel()
	bus := &fakeBus{messages: make(chan Message, 1)}
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{name: "publish_subscribe", run: func(t *testing.T) { RunPublishSubscribe(t, bus, bus) }},
		{name: "request_reply", run: func(t *testing.T) { RunRequestReply(t, bus) }},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) { t.Parallel(); tc.run(t) })
	}
}

type fakeBus struct{ messages chan Message }

func (b *fakeBus) Publish(_ context.Context, msg Message) error            { b.messages <- msg; return nil }
func (b *fakeBus) Subscribe(context.Context, string) (Subscription, error) { return b, nil }
func (b *fakeBus) Receive(context.Context) (Message, error)                { return <-b.messages, nil }
func (b *fakeBus) Close() error                                            { return nil }
func (b *fakeBus) Request(_ context.Context, msg Message) (Message, error) {
	msg.Value = []byte("pong")
	return msg, nil
}
