package repo

import (
	"fmt"
	"strings"

	"github.com/FourSigma/alertd/pkg/util"
)

func NewStmtGenerator(efs util.FieldSet, kfs util.FieldSet) genStmt {
	fls, _, _ := efs.Args()
	kls, _, _ := kfs.Args()
	table := efs.Name()

	return genStmt{
		table:      table,
		fls:        fls,
		kfls:       kfls,
		modifierFn: util.CamelCaseToUnderscore,
	}
}

type genStmt struct {
	table string

	fls  []string
	kfls []string
	mFn  func(string) string
}

func (g genStmt) GenericCreateStmt() (stmt string) {
	stmt = fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s) RETURNING *", g.modifierFn(g.table), strings.Join(Modify(g.fls, g.mFn), ", "), Placeholder(len(fls)))
	fmt.Println(stmt)
	return
}

func (g genStmt) GenericDeleteStmt() (stmt string) {
	stmt = fmt.Sprintf("DELETE FROM %s WHERE (%s) IN (%s)", g.modifierFn(g.table), strings.Join(Modify(g.fls, g.mFn), ","), Placeholder(len(kfls)))
	fmt.Println(stmt)
	return
}

func (g genStmt) GenericGetStmt() (stmt string) {
	stmt = fmt.Sprintf(
		"SELECT * FROM %s WHERE (%s) IN (%s)",
		g.modifierFn(g.table),
		strings.Join(Modify(g.fls, g.mFn), ","),
		Placeholder(len(fls)),
	)
	fmt.Println(stmt, args)
	return
}

func ModifyStringList(ls []string, mFn func(string) string) (rs []string) {
	rs = make([]string, len(ls))
	for i, _ := range ls {
		rs[i] = mFn(ls[i])
	}
	return rs
}

func ModifyString(src string, mFn func(string) string) string {
	return mFn(src)
}
