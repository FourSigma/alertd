package sqlhelpers

import (
	"fmt"
	"strings"

	"github.com/FourSigma/alertd/pkg/util"
)

func NewStmtGenerator(efs util.FieldSet, kfs util.FieldSet) *StmtGenerator {
	fls, _, _ := efs.Args()
	kfls, _, _ := kfs.Args()
	table := efs.Name()

	return &StmtGenerator{
		table:       table,
		fls:         fls,
		kfls:        kfls,
		mFn:         util.CamelCaseToUnderscore,
		placeHolder: PostgresPlaceholder,
	}
}

type StmtGenerator struct {
	table string

	fls         []string            // All Fields
	kfls        []string            // Key Fields
	mFn         func(string) string //String modifier function (ex. CamelCase -> camel_case)
	placeHolder func(int) string    // Database placeholders $1

	cache struct {
		insertStmt string
		deleteStmt string
		selectStmt string
		getStmt    string
		updateStmt map[string]string
	}
}

const (
	tmplInsertStmt = "INSERT INTO %s(%s) VALUES (%s) RETURNING *"
	tmplDeleteStmt = "DELETE FROM %s WHERE (%s) IN (%s)"
	tmplGetStmt    = "SELECT * FROM %s WHERE (%s) IN (%s)"
	tmplSelectStmt = "SELECT * FROM %s"
)

func (g *StmtGenerator) genAttributeStmt(tmpl string) string {
	return fmt.Sprintf(tmpl, g.mFn(g.table), strings.Join(ModifyStringList(g.fls, g.mFn), ", "), g.placeHolder(len(g.fls)))
}

func (g *StmtGenerator) genKeyStmt(tmpl string) string {
	return fmt.Sprintf(tmpl, g.mFn(g.table), strings.Join(ModifyStringList(g.kfls, g.mFn), ", "), g.placeHolder(len(g.kfls)))
}

func (g *StmtGenerator) InsertStmt() string {
	if g.cache.insertStmt != "" {
		return g.cache.insertStmt
	}
	g.cache.insertStmt = g.genAttributeStmt(tmplInsertStmt)
	return g.cache.insertStmt
}

func (g *StmtGenerator) DeleteStmt() string {
	if g.cache.deleteStmt != "" {
		return g.cache.deleteStmt
	}
	g.cache.deleteStmt = g.genKeyStmt(tmplDeleteStmt)
	return g.cache.deleteStmt
}

func (g *StmtGenerator) GetStmt() string {
	if g.cache.getStmt != "" {
		return g.cache.getStmt
	}
	g.cache.getStmt = g.genKeyStmt(tmplGetStmt)
	return g.cache.getStmt
}

func (g *StmtGenerator) SelectStmt() string {
	if g.cache.selectStmt != "" {
		return g.cache.selectStmt
	}
	g.cache.selectStmt = fmt.Sprintf(tmplSelectStmt, g.mFn(g.table))
	return g.cache.selectStmt
}

func (g *StmtGenerator) UpdateStmt(dfn []string) string {
	hash := strings.Join(dfn, ":")
	if val, ok := g.cache.updateStmt[hash]; ok {
		return val
	}
	g.cache.updateStmt[hash] = BuildUpdateQuery(g.table, dfn, g.kfls)
	return g.cache.updateStmt[hash]
}
