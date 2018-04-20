package sqlhelpers

import (
	"context"
	"fmt"

	"github.com/FourSigma/alertd/pkg/util"
)

func Insert(ctx context.Context, gen *StmtGenerator, fs util.FieldSet) (err error) {
	db, err := GetQueryerFromContext(ctx)
	if err != nil {
		return err
	}
	fmt.Println(gen.InsertStmt())
	if err = db.QueryRowxContext(ctx, gen.InsertStmt(), fs.Vals()...).Scan(fs.Ptrs()...); err != nil {
		return
	}
	return
}

func Get(ctx context.Context, gen *StmtGenerator, key util.FieldSet, dest util.FieldSet) (err error) {
	db, err := GetQueryerFromContext(ctx)
	if err != nil {
		return err
	}

	fmt.Println(gen.GetStmt())
	if err = db.QueryRowxContext(ctx, gen.GetStmt(), key.Vals()...).Scan(dest.Ptrs()...); err != nil {
		return
	}
	return
}

func Delete(ctx context.Context, gen *StmtGenerator, key util.FieldSet) (err error) {
	db, err := GetQueryerFromContext(ctx)
	if err != nil {
		return err
	}
	fmt.Println(gen.DeleteStmt())
	if _, err = db.ExecContext(ctx, gen.DeleteStmt(), key.Vals()...); err != nil {
		return
	}
	return
}

func Update(ctx context.Context, gen *StmtGenerator, kFS util.FieldSet, dbFS util.FieldSet, mFS util.FieldSet) (isEmpty bool, err error) {
	dfn, targs, isEmpty := UpdateFieldSetDiff(mFS, dbFS, kFS)
	if isEmpty {
		isEmpty = true
		return
	}

	db, err := GetQueryerFromContext(ctx)
	if err != nil {
		return
	}

	stmt := gen.UpdateStmt(dfn)
	fmt.Println(stmt)
	if err = db.QueryRowxContext(ctx, stmt, targs...).Scan(mFS.Ptrs()...); err != nil {
		return
	}
	return
}

func Select(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	db, err := GetQueryerFromContext(ctx)
	if err != nil {
		return err
	}

	fmt.Println(query)
	if err = db.SelectContext(ctx, dest, query, args...); err != nil {
		return
	}
	return
}
