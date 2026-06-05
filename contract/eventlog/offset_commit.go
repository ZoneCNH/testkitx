package eventlog

import (
	"context"

	"github.com/ZoneCNH/testkitx/requirex"
)

type OffsetCommitter interface {
	CommitOffset(context.Context, string, int64) error
}

func RunOffsetCommit(t requirex.TestingT, committer OffsetCommitter) {
	t.Helper()
	requirex.NoError(t, committer.CommitOffset(context.Background(), "testkitx.contract.eventlog", 1))
}
