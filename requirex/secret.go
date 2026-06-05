package requirex

import (
	"fmt"
	"strings"
)

func NoSecretLeak(t TestingT, value any, secrets ...string) {
	t.Helper()
	rendered := fmt.Sprintf("%+v", value)
	for _, secret := range secrets {
		if secret == "" {
			continue
		}
		if strings.Contains(rendered, secret) {
			t.Fatalf("secret value leaked in output")
		}
	}
}
