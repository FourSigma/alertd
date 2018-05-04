package sqlhelpers

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/FourSigma/alertd/pkg/util"
)

func NewStmtGenerator(schema string, entity util.Entity, key util.EntityKey, pural bool) StmtGenerator {
	fls, _, _ := entity.FieldSet().Args()
	kfls, _, _ := key.FieldSet().Args()
	table := entity.FieldSet().Name()

	//Convert from CamelCase to under_score
	fls = ModifyStringList(fls, util.CamelCaseToUnderscore)
	kfls = ModifyStringList(kfls, util.CamelCaseToUnderscore)
	table = ModifyString(table, util.CamelCaseToUnderscore)

	if schema != "" {
		table = schema + "." + table
	}
	if pural {
		table = table + "s"
	}
	return StmtGenerator{
		table:       table,
		fls:         fls,
		kfls:        kfls,
		placeHolder: PostgresPlaceholder,
	}
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
	tmplGetStmt    = "SELECT %s FROM %s WHERE (%s) IN (%s)"
	tmplSelectStmt = "SELECT * FROM %s"
)

func (g StmtGenerator) genAttributeStmt(tmpl string) *bytes.Buffer {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, tmpl, g.table, g.JoinedColumnFields(), g.placeHolder(g.AttributeLen()))
	return buf
}

func (g StmtGenerator) genGetStmt() *bytes.Buffer {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, tmplGetStmt, g.JoinedColumnFields(), g.table, g.JoinedKeyFields(), g.placeHolder(g.KeyLen()))
	return buf
}

func (g StmtGenerator) genKeyStmt(tmpl string) *bytes.Buffer {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, tmpl, g.table, g.JoinedKeyFields(), g.placeHolder(g.KeyLen()))
	return buf
}

func (g StmtGenerator) KeyLen() int {
	return len(g.kfls)
}
func (g StmtGenerator) AttributeLen() int {
	return len(g.fls)
}
func (g StmtGenerator) KeyFieldNames() []string {
	return g.kfls
}

func (g StmtGenerator) ColumnFieldNames() []string {
	return g.fls
}

func (g StmtGenerator) JoinedKeyFields() string {
	return strings.Join(g.KeyFieldNames(), ", ")
}

func (g StmtGenerator) JoinedColumnFields() string {
	return strings.Join(g.ColumnFieldNames(), ", ")
}

func (g StmtGenerator) InsertStmt() string {
	return g.genAttributeStmt(tmplInsertStmt).String()
}

func (g StmtGenerator) DeleteStmt() string {
	return g.genKeyStmt(tmplDeleteStmt).String()
}

func (g StmtGenerator) GetStmt() string {
	return g.genGetStmt().String()
}

func (g StmtGenerator) SelectStmt() *bytes.Buffer {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, tmplSelectStmt, g.table)
	return buf
}

func (g StmtGenerator) UpdateStmt(dfn []string) string {
	dfn = ModifyStringList(dfn, util.CamelCaseToUnderscore)
	return BuildUpdateQuery(g.table, dfn, g.kfls).String()
}
