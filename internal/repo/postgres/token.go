package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/FourSigma/alertd/internal/core"
	"github.com/FourSigma/alertd/pkg/sqlhelpers"
	_ "github.com/lib/pq"
)

type tokenRepo struct {
	crud sqlhelpers.CRUD
}

func (u tokenRepo) Create(ctx context.Context, token *core.Token) (err error) {
	token.CreatedAt = time.Now()
	token.UpdatedAt = time.Now()
	return u.crud.Insert(ctx, token)
}

func (u tokenRepo) Get(ctx context.Context, key core.TokenKey) (tkn *core.Token, err error) {
	tkn = &core.Token{}
	if err = u.crud.Get(ctx, key, tkn); err != nil {
		return
	}
	return
}

func (u tokenRepo) Delete(ctx context.Context, key core.TokenKey) (err error) {
	return u.crud.Delete(ctx, key)
}

func (u tokenRepo) List(ctx context.Context, filt core.TokenFilter, opts ...core.Opts) (ls core.TokenList, err error) {
	if err = filt.Valid(); err != nil {
		return
	}

	var args []interface{}
	qbuf := u.crud.StmtGenerator().SelectStmt()

	switch typ := filt.(type) {
	case core.FilterTokenAll:
	case core.FilterTokenByStateId:
		fmt.Fprintf(qbuf, " WHERE state_id = '%s'", string(typ.StateId))

	case *core.FilterTokenKeyIn:
		kls, total, keyLen := typ.KeyList, len(typ.KeyList), u.crud.StmtGenerator().KeyLen()
		args = make([]interface{}, total*keyLen)
		for i, j := 0, 0; i < len(kls); i, j = i+1, j+keyLen {
			k := kls[i].FieldSet().Vals()
			args[j], args[j+1] = k[0], k[1]
		}
		fmt.Fprintf(qbuf, " WHERE (id) IN %s ", sqlhelpers.PlaceholderKeyIn(total, keyLen))

	case *core.FilterTokenUserKeyIn:
		kls, total, keyLen := typ.KeyList, len(typ.KeyList), len((core.UserKey{}).FieldSet().Vals())
		args = make([]interface{}, total*keyLen)

		for i, j := 0, 0; i < len(kls); i, j = i+1, j+keyLen {
			k := kls[i].FieldSet().Vals()
			args[j], args[j+1] = k[0], k[1]
		}
		fmt.Fprintf(qbuf, " WHERE (id) IN %s ", sqlhelpers.PlaceholderKeyIn(total, keyLen))

	default:
		err = fmt.Errorf("Unknown TokenFilter Type %#v", typ)
		return
	}

	if err = u.crud.Select(ctx, &ls, qbuf.String(), args); err != nil {
		return
	}
	return
}

func (u tokenRepo) Update(ctx context.Context, key core.TokenKey, tkn *core.Token) (err error) {
	//Get from database
	dbTkn, err := u.Get(ctx, key)
	if err != nil {
		return
	}

	tkn.UpdatedAt = time.Now()

	isEmpty, err := u.crud.Update(ctx, key, dbTkn, tkn)
	if err != nil {
		return
	}
	if isEmpty {
		*tkn = *dbTkn
		return
	}
	return
}
