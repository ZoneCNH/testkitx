package columnstore

import (
	"context"

	"github.com/ZoneCNH/testkitx/requirex"
)

type BatchInserter interface {
	InsertBatch(context.Context, string, []Row) (Result, error)
}

func RunBatchInsert(t requirex.TestingT, inserter BatchInserter) {
	t.Helper()
	result, err := inserter.InsertBatch(context.Background(), "testkitx_contract_columnstore", []Row{{"id": "1"}, {"id": "2"}})
	requirex.NoError(t, err)
	if result.Rows < 2 {
		t.Fatalf("expected at least 2 inserted rows, got %d", result.Rows)
	}
}
