package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/FourSigma/alertd/internal/core"
	"github.com/FourSigma/alertd/pkg/sqlhelpers"
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

type userRepo struct {
	gen *sqlhelpers.StmtGenerator
}

func (u userRepo) Create(ctx context.Context, user *core.User) (err error) {
	user.Id = uuid.NewV4()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
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
	qbuf := u.gen.SelectStmt()

	switch typ := filt.(type) {

	case core.FilterUserAll:

	case core.FilterUserByStateId:
		fmt.Fprintf(qbuf, " WHERE state_id = '%s'", string(typ.StateId))

	case *core.FilterUserKeyIn:
		kls, total, keyLen := typ.KeyList, len(typ.KeyList), u.gen.KeyLen()
		args = make([]interface{}, total*keyLen)

		for i, j := 0, 0; i < len(kls); i, j = i+1, j+keyLen {
			k := kls[i].FieldSet().Vals()
			args[j] = k[0]
		}
		fmt.Fprintf(qbuf, " WHERE (id) IN %s ", sqlhelpers.PlaceholderKeyIn(total, keyLen))

	default:
		err = fmt.Errorf("Unknown UserFilter Type %#v", typ)
		return
	}

	if err = sqlhelpers.Select(ctx, &ls, qbuf.String(), args...); err != nil {
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
