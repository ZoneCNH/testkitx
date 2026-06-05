package common

import (
	"fmt"

	"github.com/ZoneCNH/testkitx/requirex"
)

const ErrorStandardKindID = "common.error.standard_error_kind"

type Error struct {
	Kind    string
	Op      string
	Message string
}

func (e Error) Error() string {
	if e.Op == "" {
		return fmt.Sprintf("%s: %s", e.Kind, e.Message)
	}
	return fmt.Sprintf("%s: %s: %s", e.Kind, e.Op, e.Message)
}

func (e Error) ErrorKind() string { return e.Kind }

func RunStandardErrorKind(t requirex.TestingT, err error, want string) {
	t.Helper()
	requirex.ErrorKind(t, err, want)
}
