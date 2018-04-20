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
	user.CreatedAt = time.Now()
	return sqlhelpers.Insert(ctx, u.gen, user.FieldSet())
}

func (u userRepo) Get(ctx context.Context, key core.UserKey) (usr *core.User, err error) {
	usr = &core.User{}
	if err = sqlhelpers.Get(ctx, u.gen, key.FieldSet(), usr.FieldSet()); err != nil {
		return
	}
	return
}

func (u userRepo) Delete(ctx context.Context, key core.UserKey) (err error) {
	return sqlhelpers.Delete(ctx, u.gen, key.FieldSet())
}

func (u userRepo) List(ctx context.Context, filt core.UserFilter, opts ...core.Opts) (ls core.UserList, err error) {
	if err = filt.Valid(); err != nil {
		return
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
			//Need to refactor composite primary keys (more than one key)
			args[i] = v.FieldSet().Vals()[0]
		}

	default:
		err = fmt.Errorf("Unknown UserFilter Type %#v", typ)
		return
	}

	if err = sqlhelpers.Select(ctx, &ls, query, args...); err != nil {
		return
	}
	return
}
func (u userRepo) Update(ctx context.Context, key core.UserKey, usr *core.User) (err error) {
	//Get from database
	dbUsr, err := u.Get(ctx, key)
	if err != nil {
		return
	}

	usr.UpdatedAt = time.Now()

	isEmpty, err := sqlhelpers.Update(ctx, u.gen, key.FieldSet(), dbUsr.FieldSet(), usr.FieldSet())
	if err != nil {
		return
	}
	if isEmpty {
		*usr = *dbUsr
		return
	}

	return
}
