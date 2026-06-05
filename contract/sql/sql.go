package sql

import (
	"context"

	"github.com/ZoneCNH/testkitx/requirex"
)

type Result interface{ RowsAffected() int64 }

type Rows interface {
	Next() bool
	Close() error
}

type DB interface {
	Exec(context.Context, string, ...any) (Result, error)
	Query(context.Context, string, ...any) (Rows, error)
}

func RunExecQuery(t requirex.TestingT, db DB) {
	t.Helper()
	ctx := context.Background()
	_, err := db.Exec(ctx, "testkitx contract exec")
	requirex.NoError(t, err)
	rows, err := db.Query(ctx, "testkitx contract query")
	requirex.NoError(t, err)
	if rows == nil {
		t.Fatalf("expected non-nil rows")
	}
	requirex.NoError(t, rows.Close())
}
