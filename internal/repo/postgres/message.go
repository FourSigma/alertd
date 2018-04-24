package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/FourSigma/alertd/internal/core"
	"github.com/FourSigma/alertd/pkg/sqlhelpers"
	uuid "github.com/satori/go.uuid"
)

type messageRepo struct {
	crud sqlhelpers.CRUD
}

func (u messageRepo) Create(ctx context.Context, message *core.Message) (err error) {
	message.Id = uuid.NewV4()
	message.CreatedAt = time.Now()
	message.UpdatedAt = time.Now()
	return u.crud.Insert(ctx, message)
}

func (u messageRepo) Get(ctx context.Context, key core.MessageKey) (msg *core.Message, err error) {
	msg = &core.Message{}
	if err = u.crud.Get(ctx, key, msg); err != nil {
		return
	}
	return
}

func (u messageRepo) Delete(ctx context.Context, key core.MessageKey) (err error) {
	return u.crud.Delete(ctx, key)
}

func (u messageRepo) List(ctx context.Context, filt core.MessageFilter, opts ...core.Opts) (ls core.MessageList, err error) {
	if err = filt.Valid(); err != nil {
		return
	}

	var args []interface{}
	qbuf := u.crud.StmtGenerator().SelectStmt()

	switch typ := filt.(type) {

	case core.FilterMessageAll:

	case core.FilterMessageByTypeId:
		fmt.Fprintf(qbuf, " WHERE type_id = '%s'", string(typ.TypeId))

	case *core.FilterMessageKeyIn:
		kls, total, keyLen := typ.KeyList, len(typ.KeyList), u.crud.StmtGenerator().KeyLen()
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

	if err = u.crud.Select(ctx, &ls, qbuf.String(), args); err != nil {
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

	isEmpty, err := u.crud.Update(ctx, key, dbMsg, msg)
	if err != nil {
		return
	}
	if isEmpty {
		*msg = *dbMsg
		return
	}

	return
}
