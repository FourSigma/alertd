package sqlhelpers

import (
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
	table string

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

func (g *StmtGenerator) genAttributeStmt(tmpl string) string {
	//Original revert back
	return fmt.Sprintf(tmpl, g.table, strings.Join(g.fls, ", "), g.placeHolder(len(g.fls)))
}

func (g *StmtGenerator) genKeyStmt(tmpl string) string {
	return fmt.Sprintf(tmpl, g.table, strings.Join(g.kfls, ", "), g.placeHolder(len(g.kfls)))
}

func (g *StmtGenerator) InsertStmt() string {
	return g.genAttributeStmt(tmplInsertStmt)
}

func (g *StmtGenerator) DeleteStmt() string {
	return g.genKeyStmt(tmplDeleteStmt)
}

func (g *StmtGenerator) GetStmt() string {
	return g.genKeyStmt(tmplGetStmt)
}

func (g *StmtGenerator) SelectStmt() string {
	return fmt.Sprintf(tmplSelectStmt, g.table)
}

func (g *StmtGenerator) UpdateStmt(dfn []string) string {
	dfn = ModifyStringList(dfn, util.CamelCaseToUnderscore)
	return BuildUpdateQuery(g.table, dfn, g.kfls)
}
