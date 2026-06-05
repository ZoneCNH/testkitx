package pubsub

import (
	"context"

	"github.com/ZoneCNH/testkitx/requirex"
)

type RequestReplier interface {
	Request(context.Context, Message) (Message, error)
}

func RunRequestReply(t requirex.TestingT, rr RequestReplier) {
	t.Helper()
	msg, err := rr.Request(context.Background(), Message{Topic: "testkitx.contract.request", Value: []byte("ping")})
	requirex.NoError(t, err)
	if len(msg.Value) == 0 {
		t.Fatalf("expected reply payload")
	}
}
