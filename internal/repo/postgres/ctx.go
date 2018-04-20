package postgres

import (
	"context"

	"github.com/FourSigma/alertd/pkg/sqlhelpers"
)

func Transact(ctx context.Context, txFunc func(ctx context.Context) error) (err error) {
	db, err := sqlhelpers.GetSqlDbFromContext(ctx)
	if err != nil {
		return
	}
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return
	}
	ctx = context.WithValue(ctx, sqlhelpers.CtxDbKey, tx)
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	err = txFunc(ctx)
	return err
}
