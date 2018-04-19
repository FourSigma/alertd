package postgres

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

func GetDBFromContext(ctx context.Context) (db *sqlx.DB, err error) {
	var ok bool
	db, ok = ctx.Value(CtxDbKey).(*sqlx.DB)
	if !ok || db == nil {
		err = errors.New("Context Error: Database key not loaded in ctx")
		return
	}
	return
}
