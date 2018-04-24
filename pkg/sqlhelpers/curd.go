package sqlhelpers

import (
	"context"
	"fmt"

	"github.com/FourSigma/alertd/pkg/util"
)

func NewCRUD(g StmtGenerator, hErr func(error) error) CRUD {
	return CRUD{
		gen:       g,
		handleErr: hErr,
	}
}

type CRUD struct {
	gen       StmtGenerator
	handleErr func(error) error
}

func (c CRUD) StmtGenerator() StmtGenerator {
	return c.gen
}
func (c CRUD) Insert(ctx context.Context, fs util.Entity) (err error) {
	db, err := GetQueryerFromContext(ctx)
	if err != nil {
		return c.handleErr(err)
	}
	fmt.Println(c.gen.InsertStmt())
	if err = db.QueryRowxContext(ctx, c.gen.InsertStmt(), fs.FieldSet().Vals()...).Scan(fs.FieldSet().Ptrs()...); err != nil {
		return c.handleErr(err)
	}
	return
}

func (c CRUD) Get(ctx context.Context, key util.EntityKey, dest util.Entity) (err error) {
	db, err := GetQueryerFromContext(ctx)
	if err != nil {
		return c.handleErr(err)
	}

	fmt.Println(c.gen.GetStmt())
	if err = db.QueryRowxContext(ctx, c.gen.GetStmt(), key.FieldSet().Vals()...).Scan(dest.FieldSet().Ptrs()...); err != nil {
		return c.handleErr(err)
	}
	return
}

func (c CRUD) Delete(ctx context.Context, key util.EntityKey) (err error) {
	db, err := GetQueryerFromContext(ctx)
	if err != nil {
		return c.handleErr(err)
	}
	fmt.Println(c.gen.DeleteStmt())
	if _, err = db.ExecContext(ctx, c.gen.DeleteStmt(), key.FieldSet().Vals()...); err != nil {
		return c.handleErr(err)
	}
	return
}

func (c CRUD) Update(ctx context.Context, key util.EntityKey, mod util.Entity) (err error) {
	db, err := GetQueryerFromContext(ctx)
	if err != nil {
		return c.handleErr(err)
	}

	//Get entity from database for comparison
	dbEntity := mod.New()
	if err = c.Get(ctx, key, dbEntity); err != nil {
		return
	}

	dfn, targs, isEmpty := UpdateFieldSetDiff(mod.FieldSet(), dbEntity.FieldSet(), key.FieldSet())
	if isEmpty {
		mod = dbEntity
		//Assign NotModfiedError here
		return
	}

	stmt := c.gen.UpdateStmt(dfn)
	fmt.Println(stmt)
	if err = db.QueryRowxContext(ctx, stmt, targs...).Scan(mod.FieldSet().Ptrs()...); err != nil {
		return c.handleErr(err)
	}
	return
}

func (c CRUD) Select(ctx context.Context, dest interface{}, query string, args []interface{}) (err error) {
	db, err := GetQueryerFromContext(ctx)
	if err != nil {
		return c.handleErr(err)
	}

	fmt.Println(query)
	if err = db.SelectContext(ctx, dest, query, args...); err != nil {
		return c.handleErr(err)
	}
	return
}
