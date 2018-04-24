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
	crud sqlhelpers.CRUD
}

func (u userRepo) Create(ctx context.Context, user *core.User) (err error) {
	user.Id = uuid.NewV4()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return u.crud.Insert(ctx, user)
}

func (u userRepo) Get(ctx context.Context, key core.UserKey) (usr core.User, err error) {
	if err = u.crud.Get(ctx, key, &usr); err != nil {
		return
	}
	return
}

func (u userRepo) Delete(ctx context.Context, key core.UserKey) (err error) {
	return u.crud.Delete(ctx, key)
}

func (u userRepo) List(ctx context.Context, filt core.UserFilter, opts ...core.Opts) (ls core.UserList, err error) {
	if err = filt.Valid(); err != nil {
		return
	}

	var args []interface{}
	qbuf := u.crud.StmtGenerator().SelectStmt()

	switch typ := filt.(type) {

	case core.FilterUserAll:

	case core.FilterUserByStateId:
		fmt.Fprintf(qbuf, " WHERE state_id = '%s'", string(typ.StateId))

	case *core.FilterUserKeyIn:
		kls, total, keyLen := typ.KeyList, len(typ.KeyList), u.crud.StmtGenerator().KeyLen()
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

	if err = u.crud.Select(ctx, &ls, qbuf.String(), args); err != nil {
		fmt.Println(err)
		return
	}
	return
}
func (u userRepo) Update(ctx context.Context, key core.UserKey, usr *core.User) (err error) {
	usr.UpdatedAt = time.Now()
	if err = u.crud.Update(ctx, key, usr); err != nil {
		fmt.Println(err)
		return
	}
	return
}
