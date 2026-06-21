package pubsub

import (
	"context"
	"runtime"
	"sync"
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

func TestRunPublishSubscribeNilSub(t *testing.T) {
	t.Parallel()
	probe := &fatalProbeT{}
	bus := &fakeBus{messages: make(chan Message, 1)}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() { _ = recover() }()
		RunPublishSubscribe(probe, bus, nilSubBus{})
	}()
	wg.Wait()
	if !probe.failed {
		t.Fatal("expected failure for nil subscription")
	}
}

func TestRunRequestReplyEmptyValue(t *testing.T) {
	t.Parallel()
	probe := &fatalProbeT{}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() { _ = recover() }()
		RunRequestReply(probe, emptyReplyBus{})
	}()
	wg.Wait()
	if !probe.failed {
		t.Fatal("expected failure for empty reply")
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

type nilSubBus struct{}

func (nilSubBus) Subscribe(context.Context, string) (Subscription, error) { return nil, nil }
func (nilSubBus) Publish(_ context.Context, _ Message) error              { return nil }

type emptyReplyBus struct{}

func (emptyReplyBus) Request(_ context.Context, msg Message) (Message, error) {
	msg.Value = nil
	return msg, nil
}

type fatalProbeT struct{ failed bool }

func (p *fatalProbeT) Helper()               {}
func (p *fatalProbeT) Fatalf(string, ...any) { p.failed = true; runtime.Goexit() }
