package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ZoneCNH/testkitx/pkg/testkitx"
)

func main() {
	client, err := testkitx.New(context.Background(), testkitx.Config{Name: "testkitx"})
	if err != nil {
		fmt.Fprintf(os.Stderr, "create client: %v\n", err)
		return
	}

	status := client.HealthCheck(context.Background())
	fmt.Println(status.Status)
}
