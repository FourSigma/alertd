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

type StatementGenerator struct {
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

const (
	tmplInsertStmt = "INSERT INTO %s(%s) VALUES (%s) RETURNING *"
	tmplDeleteStmt = "DELETE FROM %s WHERE (%s) IN (%s)"
	tmplGetStmt    = "SELECT * FROM %s WHERE (%s) IN (%s)"
	tmplGetStmt    = "UPDATE  FROM %s WHERE (%s) IN (%s)"
)

func (g *StatementGenerator) generateStmt(tmpl string) string {
	return fmt.Sprintf(tmpl, g.modifierFn(g.table), strings.Join(Modify(g.fls, g.mFn), ", "), Placeholder(len(fls)))
}

func (g *StatementGenerator) InsertStmt() (stmt string) {
	if g.cache.insertStmt != "" {
		return g.cache.insertStmt
	}
	g.cache.insertStmt = g.generateStmt(tmplInsertStmt)
	return g.cache.insertStmt
}

func (g *StatementGenerator) DeleteStmt() (stmt string) {
	if g.cache.deleteStmt != "" {
		return g.cache.deleteStmt
	}
	g.cache.deleteStmt = g.generateStmt(tmplDeleteStmt)
	return g.cache.deleteStmt
}

func (g *StatementGenerator) GetStmt() (stmt string) {
	if g.cache.getStmt != "" {
		return g.cache.getStmt
	}
	g.cache.getStmt = g.generateStmt(tmplGetStmt)
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
