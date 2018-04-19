package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/FourSigma/alertd/internal/core"
	"github.com/FourSigma/alertd/pkg/sqlhelpers"
)

type userRepo struct {
	gen *sqlhelpers.StmtGenerator
}

func (u userRepo) Create(ctx context.Context, user *core.User) (err error) {
	db, err := GetDBFromContext(ctx)
	if err != nil {
		return err
	}
	fs := user.FieldSet()
	fmt.Println(u.gen.InsertStmt())
	if err = db.QueryRowxContext(ctx, u.gen.InsertStmt(), fs.Vals()...).Scan(fs.Ptrs()...); err != nil {
		return
	}
	return
}

func (u userRepo) Get(ctx context.Context, key core.UserKey) (usr *core.User, err error) {
	db, err := GetDBFromContext(ctx)
	if err != nil {
		return nil, err
	}

	fs := key.FieldSet()
	usr = &core.User{}
	fmt.Println(u.gen.GetStmt(), fs.Vals())
	if err = db.QueryRowxContext(ctx, u.gen.GetStmt(), fs.Vals()...).Scan(usr.FieldSet().Ptrs()...); err != nil {
		return
	}
	return
}

func (u userRepo) Delete(ctx context.Context, key core.UserKey) (err error) {
	db, err := GetDBFromContext(ctx)
	if err != nil {
		return err
	}
	fs := key.FieldSet()
	fmt.Println(u.gen.DeleteStmt())
	if _, err = db.ExecContext(ctx, u.gen.DeleteStmt(), fs.Vals()...); err != nil {
		return
	}
	return
}

func (u userRepo) List(ctx context.Context, filt core.UserFilter, opts ...core.Opts) (ls core.UserList, err error) {
	if err = filt.Valid(); err != nil {
		return
	}

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

	case *core.FilterUserKeyIn:
		total, keyLen := len(typ.KeyList), len((core.UserKey{}).FieldSet().Vals())
		query = fmt.Sprintf("%s WHERE (id) IN %s", query, sqlhelpers.PlaceholderKeyIn(total, keyLen))
		args = make([]interface{}, total*keyLen)
		for i, v := range typ.KeyList {
			args[i] = v.FieldSet().Vals()[0]
		}

	default:
		err = fmt.Errorf("Unknown UserFilter Type %#v", typ)
		return
	}

	fmt.Println(query)
	if err = db.SelectContext(ctx, &ls, query, args...); err != nil {
		return
	}

	return
}

func (u userRepo) Update(ctx context.Context, key core.UserKey, usr *core.User) (err error) {
	db, err := GetDBFromContext(ctx)
	if err != nil {
		return
	}

	//Get from database
	dbUsr, err := u.Get(ctx, key)
	if err != nil {
		return
	}

	usr.UpdatedAt = time.Now()

	mFS, dbFS, kFS := usr.FieldSet(), dbUsr.FieldSet(), key.FieldSet()
	dfn, targs, isEmpty := sqlhelpers.UpdateFieldSetDiff(mFS, dbFS, kFS)
	if isEmpty {
		*usr = *dbUsr
		return
	}

	stmt := u.gen.UpdateStmt(dfn)
	fmt.Println(stmt)
	if err = db.QueryRowx(stmt, targs...).Scan(mFS.Ptrs()...); err != nil {
		return
	}
	return
}
