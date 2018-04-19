package sqlhelpers

import (
	"fmt"
	"testing"

	"github.com/FourSigma/alertd/internal/core"
	"github.com/FourSigma/alertd/pkg/sqlhelpers"
)

func TestGenInsertStatement(tst *testing.T) {

	gen := sqlhelpers.NewStmtGenerator((&core.User{}).FieldSet(), (&core.UserKey{}).FieldSet())
	fmt.Println(gen.InsertStmt())
	fmt.Println(gen.DeleteStmt())
	fmt.Println(gen.GetStmt())
	fmt.Println(gen.SelectStmt())

	return
}
