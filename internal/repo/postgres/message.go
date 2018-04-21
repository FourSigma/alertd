package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/FourSigma/alertd/internal/core"
	"github.com/FourSigma/alertd/pkg/sqlhelpers"
)

type messageRepo struct {
	gen *sqlhelpers.StmtGenerator
}

func (u messageRepo) Create(ctx context.Context, message *core.Message) (err error) {
	message.CreatedAt = time.Now()
	return sqlhelpers.Insert(ctx, u.gen, message.FieldSet())
}

func (u messageRepo) Get(ctx context.Context, key core.MessageKey) (msg *core.Message, err error) {
	msg = &core.Message{}
	if err = sqlhelpers.Get(ctx, u.gen, key.FieldSet(), msg.FieldSet()); err != nil {
		return
	}
	return
}

func (u messageRepo) Delete(ctx context.Context, key core.MessageKey) (err error) {
	return sqlhelpers.Delete(ctx, u.gen, key.FieldSet())
}

func (u messageRepo) List(ctx context.Context, filt core.MessageFilter, opts ...core.Opts) (ls core.MessageList, err error) {
	if err = filt.Valid(); err != nil {
		return
	}

	var args []interface{}
	qbuf := u.gen.SelectStmt()

	switch typ := filt.(type) {

	case core.FilterMessageAll:

	case core.FilterMessageByTypeId:
		fmt.Fprintf(qbuf, " WHERE type_id = '%s'", string(typ.TypeId))

	case *core.FilterMessageKeyIn:
		kls, total, keyLen := typ.KeyList, len(typ.KeyList), u.gen.KeyLen()
		args = make([]interface{}, total*keyLen)

		for i, j := 0, 0; i < len(kls); i, j = i+1, j+keyLen {
			k := kls[i].FieldSet().Vals()
			args[j] = k[0]
		}
		fmt.Fprintf(qbuf, " WHERE (id) IN %s ", sqlhelpers.PlaceholderKeyIn(total, keyLen))

	default:
		err = fmt.Errorf("Unknown MessageFilter Type %#v", typ)
		return
	}

	if err = sqlhelpers.Select(ctx, &ls, qbuf.String(), args...); err != nil {
		return
	}
	return
}
func (u messageRepo) Update(ctx context.Context, key core.MessageKey, msg *core.Message) (err error) {
	//Get from database
	dbMsg, err := u.Get(ctx, key)
	if err != nil {
		return
	}

	msg.UpdatedAt = time.Now()

	isEmpty, err := sqlhelpers.Update(ctx, u.gen, key.FieldSet(), dbMsg.FieldSet(), msg.FieldSet())
	if err != nil {
		return
	}
	if isEmpty {
		*msg = *dbMsg
		return
	}

	return
}
