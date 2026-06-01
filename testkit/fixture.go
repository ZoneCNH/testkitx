package testkit

import (
	"time"

	"github.com/ZoneCNH/testkitx/pkg/testkitx"
)

func Config(name string) testkitx.Config {
	return testkitx.Config{
		Name:    name,
		Timeout: time.Second,
	}
}
