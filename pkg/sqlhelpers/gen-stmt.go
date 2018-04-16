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

	cache struct {
		insertStmt string
		deleteStmt string
		getStmt    string
	}
}

func (g genStmt) GenericInsertStmt() (stmt string) {
	if g.cache.insertStmt != "" {
		return g.cache.insertStmt
	}
	g.cache.insertStmt = fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s) RETURNING *", g.modifierFn(g.table), strings.Join(Modify(g.fls, g.mFn), ", "), Placeholder(len(fls)))
	return g.cache.insertStmt
}

func (g genStmt) GenericDeleteStmt() (stmt string) {
	if g.cache.deleteStmt != "" {
		return g.cache.deleteStmt
	}
	g.cache.deleteStmt = fmt.Sprintf("DELETE FROM %s WHERE (%s) IN (%s)", g.modifierFn(g.table), strings.Join(Modify(g.fls, g.mFn), ","), Placeholder(len(kfls)))
	return g.cache.deleteStmt
}

func (g genStmt) GenericGetStmt() (stmt string) {
	if g.cache.getStmt != "" {
		return g.cache.getStmt
	}
	g.cache.getStmt = fmt.Sprintf(
		"SELECT * FROM %s WHERE (%s) IN (%s)",
		g.modifierFn(g.table),
		strings.Join(Modify(g.fls, g.mFn), ","),
		Placeholder(len(fls)),
	)
	return g.cache.getStmt
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
