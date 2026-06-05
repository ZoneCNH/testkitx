package pubsub

import (
	"context"

	"github.com/ZoneCNH/testkitx/requirex"
)

type Message struct {
	Topic   string
	Key     string
	Value   []byte
	Headers map[string]string
}

type Publisher interface {
	Publish(context.Context, Message) error
}

type Subscription interface {
	Receive(context.Context) (Message, error)
	Close() error
}

type Subscriber interface {
	Subscribe(context.Context, string) (Subscription, error)
}

func RunPublishSubscribe(t requirex.TestingT, publisher Publisher, subscriber Subscriber) {
	t.Helper()
	ctx := context.Background()
	sub, err := subscriber.Subscribe(ctx, "testkitx.contract.pubsub")
	requirex.NoError(t, err)
	if sub == nil {
		t.Fatalf("expected non-nil subscription")
	}
	requirex.NoError(t, publisher.Publish(ctx, Message{Topic: "testkitx.contract.pubsub", Value: []byte("value")}))
	_, err = sub.Receive(ctx)
	requirex.NoError(t, err)
	requirex.NoError(t, sub.Close())
}
