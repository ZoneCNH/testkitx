package assertx_test

import (
	"errors"
	"testing"
	"time"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/assertx"
)

func TestEventuallySucceedsAfterRetry(t *testing.T) {
	attempts := 0
	assertx.Eventually(t, time.Second, time.Millisecond, func() error {
		attempts++
		if attempts == 3 {
			return nil
		}
		return errors.New("not yet")
	})
	assertx.Equal(t, 3, attempts)
}
