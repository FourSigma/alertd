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
	gen *sqlhelpers.StmtGenerator
}

func (u tokenRepo) Create(ctx context.Context, token *core.Token) (err error) {
	token.CreatedAt = time.Now()
	return sqlhelpers.Insert(ctx, u.gen, token.FieldSet())
}

func (u tokenRepo) Get(ctx context.Context, key core.TokenKey) (tkn *core.Token, err error) {
	tkn = &core.Token{}
	if err = sqlhelpers.Get(ctx, u.gen, key.FieldSet(), tkn.FieldSet()); err != nil {
		return
	}
	return
}

func (u tokenRepo) Delete(ctx context.Context, key core.TokenKey) (err error) {
	return sqlhelpers.Delete(ctx, u.gen, key.FieldSet())
}

func (u tokenRepo) List(ctx context.Context, filt core.TokenFilter, opts ...core.Opts) (ls core.TokenList, err error) {
	if err = filt.Valid(); err != nil {
		return
	}

	var args []interface{}
	qbuf := u.gen.SelectStmt()

	switch typ := filt.(type) {
	case core.FilterTokenAll:
	case core.FilterTokenByStateId:
		fmt.Fprintf(qbuf, " WHERE state_id = '%s'", string(typ.StateId))

	case *core.FilterTokenKeyIn:
		kls, total, keyLen := typ.KeyList, len(typ.KeyList), u.gen.KeyLen()
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

	if err = sqlhelpers.Select(ctx, &ls, qbuf.String(), args...); err != nil {
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

	isEmpty, err := sqlhelpers.Update(ctx, u.gen, key.FieldSet(), dbTkn.FieldSet(), tkn.FieldSet())
	if err != nil {
		return
	}
	if isEmpty {
		*tkn = *dbTkn
		return
	}
	return
}
