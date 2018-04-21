package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/FourSigma/alertd/internal/core"
	"github.com/FourSigma/alertd/pkg/sqlhelpers"
	_ "github.com/lib/pq"
)

type topicRepo struct {
	gen *sqlhelpers.StmtGenerator
}

func (u topicRepo) Create(ctx context.Context, topic *core.Topic) (err error) {
	topic.CreatedAt = time.Now()
	return sqlhelpers.Insert(ctx, u.gen, topic.FieldSet())
}

func (u topicRepo) Get(ctx context.Context, key core.TopicKey) (tp *core.Topic, err error) {
	tp = &core.Topic{}
	if err = sqlhelpers.Get(ctx, u.gen, key.FieldSet(), tp.FieldSet()); err != nil {
		return
	}
	return
}

func (u topicRepo) Delete(ctx context.Context, key core.TopicKey) (err error) {
	return sqlhelpers.Delete(ctx, u.gen, key.FieldSet())
}

func (u topicRepo) List(ctx context.Context, filt core.TopicFilter, opts ...core.Opts) (ls core.TopicList, err error) {
	if err = filt.Valid(); err != nil {
		return
	}

	var args []interface{}
	qbuf := u.gen.SelectStmt()

	switch typ := filt.(type) {
	case core.FilterTopicAll:

	case *core.FilterTopicKeyIn:
		kls, total, keyLen := typ.KeyList, len(typ.KeyList), u.gen.KeyLen()
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

	if err = sqlhelpers.Select(ctx, &ls, qbuf.String(), args...); err != nil {
		return
	}
	return
}

func (u topicRepo) Update(ctx context.Context, key core.TopicKey, tp *core.Topic) (err error) {
	//Get from database
	dbTp, err := u.Get(ctx, key)
	if err != nil {
		return
	}

	tp.UpdatedAt = time.Now()

	isEmpty, err := sqlhelpers.Update(ctx, u.gen, key.FieldSet(), dbTp.FieldSet(), tp.FieldSet())
	if err != nil {
		return
	}
	if isEmpty {
		*tp = *dbTp
		return
	}
	return
}
