package sqlhelpers

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

type repoPostgresCtxKey string

const (
	CtxDbKey repoPostgresCtxKey = "[Postgres] DB Key"
)

func GetQueryerFromContext(ctx context.Context) (db Queryer, err error) {
	var ok bool
	db, ok = ctx.Value(CtxDbKey).(Queryer)
	if !ok || db == nil {
		err = errors.New("Context Error: Queryer not loaded in ctx")
		return
	}
	return
}

func GetSqlDbFromContext(ctx context.Context) (db *sqlx.DB, err error) {
	var ok bool
	db, ok = ctx.Value(CtxDbKey).(*sqlx.DB)
	if !ok || db == nil {
		err = errors.New("Context Error: sqlx.DB not loaded in ctx")
		return
	}
	return
}

type Queryer interface {
	sqlx.ExtContext
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}
