package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/FourSigma/alertd/internal/core"
	"github.com/FourSigma/alertd/pkg/sqlhelpers"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func TestInsert(tst *testing.T) {
	p := userRepo{
		gen: sqlhelpers.NewStmtGenerator("alertd", (&core.User{}).FieldSet(), (&core.UserKey{}).FieldSet()),
	}
	db, err := sqlx.Connect("postgres", "user=sivawork dbname=sivawork sslmode=disable")
	if err != nil {
		tst.Error("Error connecting to database: ", err)
		return
	}
	ctx := context.WithValue(context.Background(), CtxDbKey, db)
	u := core.NewUser("TestFirstName", "TestLastName", "test@email.com", "TestPassword")
	u.PasswordHash = "TestHash"
	u.PasswordSalt = "TestSalt"

	if err = p.Create(ctx, u); err != nil {
		tst.Error(err)
		return
	}

	usr, err := p.Get(ctx, u.Key())
	if err != nil {
		tst.Error(err)
		return
	}

	// if *usr != *u {
	// 	tst.Errorf("Get doesn't match...")
	// 	return
	// }

	fmt.Println(usr)
	err = p.Update(ctx, usr.Key(), usr)
	if err != nil {
		tst.Error(err)
		return
	}

	uUsr, err := p.Get(ctx, usr.Key())
	if err != nil {
		tst.Error(err)
		return
	}

	fmt.Println(uUsr.FirstName)
	rs, err := p.List(ctx, core.FilterUserAll{})
	if err != nil {
		tst.Error(err)
		return
	}

	if err = p.Delete(ctx, u.Key()); err != nil {
		tst.Error(err)
	}
	fmt.Println(rs)
}
