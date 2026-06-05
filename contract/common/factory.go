package common

import (
	"context"

	"github.com/ZoneCNH/testkitx/requirex"
)

type Runtime interface {
	Start(context.Context) error
	Ping(context.Context) error
	Close(context.Context) error
}

type Factory interface {
	New(context.Context, Config) (Runtime, error)
}

func newRuntime(t requirex.TestingT, factory Factory, cfg Config) Runtime {
	t.Helper()
	runtime, err := factory.New(context.Background(), cfg)
	requirex.NoError(t, err)
	if runtime == nil {
		t.Fatalf("factory returned nil runtime")
	}
	return runtime
}
