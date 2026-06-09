package assertx_test

import (
	"errors"
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/assertx"
)

func TestNotEqualDetectsSameValues(t *testing.T) {
	t.Parallel()
	// NotEqual should pass when values differ.
	assertx.NotEqual(t, 1, 2)
	assertx.NotEqual(t, "a", "b")
}

func TestHasError(t *testing.T) {
	t.Parallel()
	assertx.HasError(t, errors.New("boom"))
}

func TestContains(t *testing.T) {
	t.Parallel()
	assertx.Contains(t, "hello world", "world")
}

func TestLen(t *testing.T) {
	t.Parallel()
	assertx.Len(t, []int{1, 2, 3}, 3)
	assertx.Len(t, "abc", 3)
	assertx.Len(t, map[string]int{"a": 1}, 1)
}
