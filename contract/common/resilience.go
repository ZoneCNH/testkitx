package common

import (
	"context"

	"github.com/ZoneCNH/testkitx/requirex"
)

const ResilienceCancelID = "common.resilience.context_cancel"

type Operation func(context.Context) error

func RunContextCancel(t requirex.TestingT, op Operation) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	requirex.ErrorKindOneOf(t, op(ctx), "canceled", "timeout", "unavailable")
}
