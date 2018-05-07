package sqlhelpers

import (
	"context"
	"strings"

	"github.com/FourSigma/alertd/pkg/util"
	log "github.com/Sirupsen/logrus"
)

func NewCRUD(l *log.Logger, g StmtGenerator) CRUD {
	return CRUD{
		gen: g,
		log: l.WithFields(
			log.Fields{
				"layer": "repo",
				"type":  strings.ToLower(g.table),
			},
		),
	}
}

type CRUD struct {
	gen StmtGenerator
	log *log.Entry
}

func (c CRUD) StmtGenerator() StmtGenerator {
	return c.gen
}
func (c CRUD) Insert(ctx context.Context, e util.Entity) (err error) {
	db, err := GetQueryerFromContext(ctx)
	if err != nil {
		return err
	}
	stmt := c.gen.InsertStmt()
	fs := e.FieldSet()

	c.log.WithFields(log.Fields{
		"sql": stmt,
	}).Info("SQL Statement")

	if err = db.QueryRowxContext(ctx, stmt, fs.Vals()...).Scan(fs.Ptrs()...); err != nil {
		c.log.WithFields(log.Fields{
			"values": fs.Vals(),
		}).Error(err)
		return err
	}
	c.log.WithFields(log.Fields{
		"values": fs.Vals(),
	}).Info("Successfully inserted")
	return
}

func (c CRUD) Get(ctx context.Context, key util.EntityKey, dest util.Entity) (err error) {
	db, err := GetQueryerFromContext(ctx)
	if err != nil {
		return err
	}

	stmt := c.gen.GetStmt()

	c.log.WithFields(log.Fields{
		"sql": stmt,
	}).Info("SQL Statement")

	if err = db.QueryRowxContext(ctx, stmt, key.FieldSet().Vals()...).Scan(dest.FieldSet().Ptrs()...); err != nil {
		c.log.WithFields(log.Fields{
			"values": key,
		}).Error(err)
		return err
	}
	c.log.WithFields(log.Fields{
		"values": key,
	}).Info("Obtained instance")
	return
}

func (c CRUD) Delete(ctx context.Context, key util.EntityKey) (err error) {
	db, err := GetQueryerFromContext(ctx)
	if err != nil {
		return err
	}
	stmt := c.gen.DeleteStmt()

	c.log.WithFields(log.Fields{
		"sql": stmt,
	}).Info("SQL Statement")
	if _, err = db.ExecContext(ctx, stmt, key.FieldSet().Vals()...); err != nil {
		c.log.WithFields(log.Fields{
			"values": key,
		}).Error(err)
		return err
	}
	c.log.WithFields(log.Fields{
		"values": key,
	}).Info("Successfully deleted")

	return
}

func (c CRUD) Update(ctx context.Context, key util.EntityKey, mod util.Entity) (err error, didUpdate bool) {
	db, err := GetQueryerFromContext(ctx)
	if err != nil {
		return err, false
	}

	//Get entity from database for comparison
	dbEntity := mod.New()
	if err = c.Get(ctx, key, dbEntity); err != nil {
		c.log.WithFields(log.Fields{
			"values": key,
		}).Error(err)
		return
	}

	dfn, targs, isEmpty := UpdateFieldSetDiff(mod.FieldSet(), dbEntity.FieldSet(), key.FieldSet())
	if isEmpty {
		if err = c.Get(ctx, key, mod); err != nil {
			c.log.WithFields(log.Fields{
				"values": key,
			}).Error(err)
			return
		}
		return err, false
	}

	stmt := c.gen.UpdateStmt(dfn)
	c.log.WithFields(log.Fields{
		"sql": stmt,
	}).Info("SQL Statement")
	if err = db.QueryRowxContext(ctx, stmt, targs...).Scan(mod.FieldSet().Ptrs()...); err != nil {
		c.log.WithFields(log.Fields{
			"values": targs,
		}).Error(err)
		return err, false
	}
	c.log.WithFields(log.Fields{
		"values": targs,
	}).Info("Successfully updated")

	didUpdate = true
	return
}

func (c CRUD) Select(ctx context.Context, dest interface{}, query string, args []interface{}) (err error) {
	db, err := GetQueryerFromContext(ctx)
	if err != nil {
		return err
	}

	c.log.WithFields(log.Fields{
		"sql": query,
	}).Info("sql statement")

	if err = db.SelectContext(ctx, dest, query, args...); err != nil {
		c.log.WithFields(log.Fields{
			"values": args,
		}).Error(err)
		return err
	}
	return
}
