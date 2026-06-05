package common

import (
	"context"

	"github.com/ZoneCNH/testkitx/requirex"
)

const (
	LifecycleStartID = "common.lifecycle.start"
	LifecyclePingID  = "common.lifecycle.ping"
	LifecycleCloseID = "common.lifecycle.close_idempotent"
)

var TestIDs = []string{
	LifecycleStartID,
	LifecyclePingID,
	LifecycleCloseID,
	ConfigInvalidID,
	ErrorStandardKindID,
	SecretNoLeakID,
	ResilienceCancelID,
	ObservabilityMetricsID,
}

func RunLifecycleStart(t requirex.TestingT, factory Factory) {
	t.Helper()
	runtime := newRuntime(t, factory, ValidConfig())
	requirex.NoError(t, runtime.Start(context.Background()))
	requirex.NoError(t, runtime.Close(context.Background()))
}

func RunLifecyclePing(t requirex.TestingT, factory Factory) {
	t.Helper()
	runtime := newRuntime(t, factory, ValidConfig())
	requirex.NoError(t, runtime.Start(context.Background()))
	requirex.NoError(t, runtime.Ping(context.Background()))
	requirex.NoError(t, runtime.Close(context.Background()))
}

func RunLifecycleCloseIdempotent(t requirex.TestingT, factory Factory) {
	t.Helper()
	runtime := newRuntime(t, factory, ValidConfig())
	requirex.NoError(t, runtime.Start(context.Background()))
	requirex.NoError(t, runtime.Close(context.Background()))
	requirex.NoError(t, runtime.Close(context.Background()))
}
