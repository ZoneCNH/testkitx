package common

import (
	"context"

	"github.com/ZoneCNH/testkitx/requirex"
)

const ConfigInvalidID = "common.config.invalid_config"

type Config struct {
	Name   string
	Secret string
	Values map[string]string
}

func ValidConfig() Config { return Config{Name: "testkitx-contract"} }

func InvalidConfig() Config { return Config{} }

func RunInvalidConfig(t requirex.TestingT, factory Factory) {
	t.Helper()
	runtime, err := factory.New(context.Background(), InvalidConfig())
	if runtime != nil {
		t.Fatalf("expected invalid config to return nil runtime")
	}
	requirex.ErrorKind(t, err, "validation")
}
