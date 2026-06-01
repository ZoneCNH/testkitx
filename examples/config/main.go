package main

import (
	"fmt"
	"time"

	"github.com/ZoneCNH/testkitx/pkg/testkitx"
)

func main() {
	cfg := testkitx.Config{
		Name:    "testkitx",
		Timeout: time.Second,
		Secret:  "example",
	}

	fmt.Println(cfg.Sanitize().Secret)
}
