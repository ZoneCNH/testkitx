package columnstore

import (
	"context"

	"github.com/ZoneCNH/testkitx/requirex"
)

type Row map[string]any
type Result struct{ Rows int }
type Store interface {
	Insert(context.Context, string, Row) error
	Query(context.Context, string) ([]Row, error)
}

func RunColumnStore(t requirex.TestingT, store Store) {
	t.Helper()
	ctx := context.Background()
	table := "testkitx_contract_columnstore"
	requirex.NoError(t, store.Insert(ctx, table, Row{"id": "1"}))
	rows, err := store.Query(ctx, table)
	requirex.NoError(t, err)
	if rows == nil {
		t.Fatalf("expected non-nil rows")
	}
}
