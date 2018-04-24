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

type topicRepo struct {
	crud sqlhelpers.CRUD
}

func (u topicRepo) Create(ctx context.Context, topic *core.Topic) (err error) {
	topic.Id = uuid.NewV4()
	topic.CreatedAt = time.Now()
	topic.UpdatedAt = time.Now()
	return u.crud.Insert(ctx, topic)
}

func (u topicRepo) Get(ctx context.Context, key core.TopicKey) (tp core.Topic, err error) {
	if err = u.crud.Get(ctx, key, &tp); err != nil {
		return
	}
	return
}

func (u topicRepo) Delete(ctx context.Context, key core.TopicKey) (err error) {
	return u.crud.Delete(ctx, key)
}

func (u topicRepo) List(ctx context.Context, filt core.TopicFilter, opts ...core.Opts) (ls core.TopicList, err error) {
	if err = filt.Valid(); err != nil {
		return
	}

	var args []interface{}
	qbuf := u.crud.StmtGenerator().SelectStmt()

	switch typ := filt.(type) {
	case core.FilterTopicAll:

	case *core.FilterTopicKeyIn:
		kls, total, keyLen := typ.KeyList, len(typ.KeyList), u.crud.StmtGenerator().KeyLen()
		args = make([]interface{}, total*keyLen)
		for i, j := 0, 0; i < len(kls); i, j = i+1, j+keyLen {
			k := kls[i].FieldSet().Vals()
			args[j], args[j+1] = k[0], k[1]
		}
		fmt.Fprintf(qbuf, " WHERE (id) IN %s ", sqlhelpers.PlaceholderKeyIn(total, keyLen))

	case *core.FilterTopicUserKeyIn:
		kls, total, keyLen := typ.KeyList, len(typ.KeyList), len((core.UserKey{}).FieldSet().Vals())
		args = make([]interface{}, total*keyLen)

		for i, j := 0, 0; i < len(kls); i, j = i+1, j+keyLen {
			k := kls[i].FieldSet().Vals()
			args[j], args[j+1] = k[0], k[1]
		}
		fmt.Fprintf(qbuf, " WHERE (id) IN %s ", sqlhelpers.PlaceholderKeyIn(total, keyLen))

	default:
		err = fmt.Errorf("Unknown TopicFilter Type %#v", typ)
		return
	}

	if err = u.crud.Select(ctx, &ls, qbuf.String(), args); err != nil {
		return
	}
	return
}

func (u topicRepo) Update(ctx context.Context, key core.TopicKey, tp *core.Topic) (err error) {
	tp.UpdatedAt = time.Now()
	if err = u.crud.Update(ctx, key, tp); err != nil {
		return
	}
	return
}
