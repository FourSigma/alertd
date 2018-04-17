package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/FourSigma/alertd/internal/core"
	"github.com/FourSigma/alertd/pkg/sqlhelpers"
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

type userRepo struct {
	gen *sqlhelpers.StatementGenerator
}

func (u userRepo) Create(ctx context.Context, user *core.User) (err error) {
	db, err := GetDBFromContext(ctx)
	if err != nil {
		return err
	}
	if err = db.QueryRowxContext(ctx, u.gen.InsertStmt(), user.FieldSet().Vals()...).Scan(user.FieldSet().Ptrs()...); err != nil {
		return
	}
	return
}

func (u userRepo) Get(ctx context.Context, key core.UserKey) (usr *core.User, err error) {
	db, err := GetDBFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err = db.GetContext(ctx, usr, u.gen.GetStmt(), key.FieldSet().Vals()...); err != nil {
		return
	}
	return
}

func (u userRepo) Delete(ctx context.Context, key core.UserKey) (err error) {
	db, err := GetDBFromContext(ctx)
	if err != nil {
		return err
	}
	if _, err = db.ExecContext(ctx, u.gen.DeleteStmt(), key.FieldSet().Vals()...); err != nil {
		return
	}
	return
}

func (u userRepo) List(ctx context.Context, filt core.UserFilter, opts ...core.Opts) (ls []*core.User, err error) {
	db, err := GetDBFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var args []interface{}
	query := u.gen.SelectStmt()

	switch typ := filt.(type) {
	case core.FilterUserAll:
	case core.FilterUserActiveUsers:
		query = query + " WHERE state_id = 'Active'"

	case core.FilterUserKeyIn:
		total, keyLen := len(typ.KeyList), len((core.UserKey{}).Args())
		query = fmt.Sprintf("%s WHERE (id) IN %s", query, sqlhelpers.InQueryPlaceholder(total, keyLen))
		args = make([]interface{}, total*keyLen)
		for i, v := range typ.KeyList {
			args[i] = v
		}

	default:
		err = fmt.Errorf("Unknown UserFilter Type %#v", typ)
		return
	}

	if err = db.SelectContext(ctx, ls, query, args...); err != nil {
		return
	}

	return
}

func (u userRepo) Update(ctx context.Context, key core.UserKey, usr *core.User) (err error) {
	db, err := GetDBFromContext(ctx)
	if err != nil {
		return
	}

	dbUsr, err := u.Get(ctx, key)
	if err != nil {
		return
	}

	usr.UpdatedAt = time.Now()
	stmt := u.gen.UpdateStmt(usr.FieldSet(), dbUsr.FieldSet(), key.FieldSet())
	if isEmpty {
		*usr = *dbUsr
		return
	}

	if err = db.QueryRowx(stmt, targs...).Scan(usr.FieldSet().Ptr()...); err != nil {
		return
	}

	return
}
