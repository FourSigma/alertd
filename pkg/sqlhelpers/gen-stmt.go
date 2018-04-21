package sqlhelpers

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/FourSigma/alertd/pkg/util"
)

func NewStmtGenerator(schema string, efs util.FieldSet, kfs util.FieldSet) *StmtGenerator {
	fls, _, _ := efs.Args()
	kfls, _, _ := kfs.Args()
	table := efs.Name()

	//Convert from CamelCase to under_score
	fls = ModifyStringList(fls, util.CamelCaseToUnderscore)
	kfls = ModifyStringList(kfls, util.CamelCaseToUnderscore)
	table = ModifyString(table, util.CamelCaseToUnderscore)

	sg := &StmtGenerator{
		table:       schema + "." + table + "s",
		fls:         fls,
		kfls:        kfls,
		placeHolder: PostgresPlaceholder,
	}
	return sg
}

type StmtGenerator struct {
	table       string           //Table name
	fls         []string         // All Fields
	kfls        []string         // Key Fields
	placeHolder func(int) string // Database placeholders $1

}

const (
	tmplInsertStmt = "INSERT INTO %s(%s) VALUES (%s) RETURNING *"
	tmplDeleteStmt = "DELETE FROM %s WHERE (%s) IN (%s)"
	tmplGetStmt    = "SELECT * FROM %s WHERE (%s) IN (%s)"
	tmplSelectStmt = "SELECT * FROM %s"
)

func (g *StmtGenerator) genAttributeStmt(tmpl string) *bytes.Buffer {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, tmpl, g.table, strings.Join(g.fls, ", "), g.placeHolder(len(g.fls)))
	return buf
}

func (g *StmtGenerator) genKeyStmt(tmpl string) *bytes.Buffer {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, tmpl, g.table, strings.Join(g.kfls, ", "), g.placeHolder(len(g.kfls)))
	return buf
}

func (g *StmtGenerator) InsertStmt() string {
	return g.genAttributeStmt(tmplInsertStmt).String()
}

func (g *StmtGenerator) DeleteStmt() string {
	return g.genKeyStmt(tmplDeleteStmt).String()
}

func (g *StmtGenerator) GetStmt() string {
	return g.genKeyStmt(tmplGetStmt).String()
}

func (g *StmtGenerator) SelectStmt() *bytes.Buffer {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, tmplSelectStmt, g.table)
	return buf
}

func (g *StmtGenerator) UpdateStmt(dfn []string) string {
	dfn = ModifyStringList(dfn, util.CamelCaseToUnderscore)
	return BuildUpdateQuery(g.table, dfn, g.kfls).String()
}
