package sql

import (
	"context"

	"github.com/ZoneCNH/testkitx/requirex"
)

type Tx interface {
	Exec(context.Context, string, ...any) (Result, error)
	Commit(context.Context) error
	Rollback(context.Context) error
}

type Transactor interface {
	BeginTx(context.Context) (Tx, error)
}

func RunTransaction(t requirex.TestingT, txer Transactor) {
	t.Helper()
	tx, err := txer.BeginTx(context.Background())
	requirex.NoError(t, err)
	if tx == nil {
		t.Fatalf("expected non-nil transaction")
	}
	_, err = tx.Exec(context.Background(), "testkitx contract tx")
	requirex.NoError(t, err)
	requirex.NoError(t, tx.Commit(context.Background()))

	tx, err = txer.BeginTx(context.Background())
	requirex.NoError(t, err)
	if tx == nil {
		t.Fatalf("expected non-nil transaction")
	}
	_, err = tx.Exec(context.Background(), "testkitx contract tx rollback")
	requirex.NoError(t, err)
	requirex.NoError(t, tx.Rollback(context.Background()))
}
