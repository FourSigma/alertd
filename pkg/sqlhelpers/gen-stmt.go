package sqlhelpers

import (
	"fmt"
	"strings"

	"github.com/FourSigma/alertd/pkg/util"
)

func NewStmtGenerator(efs util.FieldSet, kfs util.FieldSet) *StatementGenerator {
	fls, _, _ := efs.Args()
	kfls, _, _ := kfs.Args()
	table := efs.Name()

	return &StatementGenerator{
		table:       table,
		fls:         fls,
		kfls:        kfls,
		mFn:         util.CamelCaseToUnderscore,
		placeHolder: PostgresPlaceholder,
	}
}

type StatementGenerator struct {
	table string

	fls         []string
	kfls        []string
	mFn         func(string) string
	placeHolder func(int) string

	cache struct {
		insertStmt string
		deleteStmt string
		selectStmt string
		getStmt    string
	}
}

const (
	tmplInsertStmt = "INSERT INTO %s(%s) VALUES (%s) RETURNING *"
	tmplDeleteStmt = "DELETE FROM %s WHERE (%s) IN (%s)"
	tmplGetStmt    = "SELECT * FROM %s WHERE (%s) IN (%s)"
	tmplSelectStmt = "SELECT * FROM %s"
	tmplUpdateStmt = "UPDATE  FROM %s WHERE (%s) IN (%s)"
)

func (g *StatementGenerator) genAttributeStmt(tmpl string) string {
	return fmt.Sprintf(tmpl, g.mFn(g.table), strings.Join(ModifyStringList(g.fls, g.mFn), ", "), g.placeHolder(len(g.fls)))
}

func (g *StatementGenerator) genKeyStmt(tmpl string) string {
	return fmt.Sprintf(tmpl, g.mFn(g.table), strings.Join(ModifyStringList(g.kfls, g.mFn), ", "), g.placeHolder(len(g.kfls)))
}

func (g *StatementGenerator) InsertStmt() (stmt string) {
	if g.cache.insertStmt != "" {
		return g.cache.insertStmt
	}
	g.cache.insertStmt = g.genAttributeStmt(tmplInsertStmt)
	return g.cache.insertStmt
}

func (g *StatementGenerator) DeleteStmt() (stmt string) {
	if g.cache.deleteStmt != "" {
		return g.cache.deleteStmt
	}
	g.cache.deleteStmt = g.genKeyStmt(tmplDeleteStmt)
	return g.cache.deleteStmt
}

func (g *StatementGenerator) GetStmt() (stmt string) {
	if g.cache.getStmt != "" {
		return g.cache.getStmt
	}
	g.cache.getStmt = g.genKeyStmt(tmplGetStmt)
	return g.cache.getStmt
}

func (g *StatementGenerator) SelectStmt() (stmt string) {
	if g.cache.selectStmt != "" {
		return g.cache.selectStmt
	}
	g.cache.selectStmt = g.genAttributeStmt(tmplSelectStmt)
	return g.cache.selectStmt
}

func (g *StatementGenerator) UpdateStmt() (stmt string) {
	return
}
