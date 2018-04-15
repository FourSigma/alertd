package postgres

import (
	"context"
	"errors"

	"github.com/FourSigma/alertd/internal/core"
	"github.com/FourSigma/alertd/pkg/sqlhelpers"
	"github.com/jmoiron/sqlx"
)

func GetDBFromContext(ctx context.Context) (db *sqlx.DB, err error) {
	var ok bool
	db, ok = ctx.Value(CtxDbKey).(*sqlx.DB)
	if !ok || db == nil {
		err = errors.New("Context Error: Database key not loaded in ctx")
		return
	}
	return
}

type userRepo struct{}

func (u userRepo) Create(ctx context.Context, user *core.User) (err error) {
	db, err := GetDBFromContext(ctx)
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx,
		`INSERT INTO
           user(id, first_name, last_name, email, password, password_salt, password_hash, state_id, created_at, updated_at)
           VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		user.Id, user.FirstName, user.LastName, user.Email, user.Password, user.PasswordSalt, user.PasswordHash, user.StateId, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return
	}
	return
}

func (u userRepo) Delete(ctx context.Context, key core.UserKey) (err error) {
	db, err := GetDBFromContext(ctx)
	if err != nil {
		return err
	}
	if _, err = db.ExecContext(ctx, `DELETE FROM user WHERE (id) IN ($1)`, key.Id); err != nil {
		return
	}
	return
}

// func (_ userRepo) List(context.Context, core.UserFilter, ...core.Opts) ([]*core.User, error) {
// 	return
// }

func (u userRepo) Get(ctx context.Context, key core.UserKey) (usr *core.User, err error) {
	db, err := GetDBFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return u.get(ctx, db, key)
}

func (_ userRepo) get(ctx context.Context, db *sqlx.DB, key core.UserKey) (u *core.User, err error) {
	if err = db.GetContext(ctx, u, `SELECT * FROM user WHERE (id) IN ($1)`, key.Id); err != nil {
		return
	}
	return
}

func (u userRepo) Update(ctx context.Context, key core.UserKey, usr *core.User) (err error) {
	db, err := GetDBFromContext(ctx)
	if err != nil {
		return
	}

	query, args, err := u.getUpdateStmt(ctx, db, key, usr)
	if err != nil {
		return
	}

	if _, err = db.ExecContext(ctx, query, args...); err != nil {
		return
	}

	return
}

func (u userRepo) getUpdateStmt(ctx context.Context, db *sqlx.DB, key core.UserKey, mUser *core.User) (query string, args []interface{}, err error) {
	var dbUser *core.User
	dbUser, err = u.get(ctx, db, key)
	if err != nil {
		return
	}

	//Add Updateable Fields
	uf := sqlhelpers.FieldValueList{}
	if mUser.FirstName != dbUser.FirstName {
		uf.AddAttributeField("first_name", mUser.FirstName)
	}

	if mUser.LastName != dbUser.LastName {
		uf.AddAttributeField("last_name", mUser.LastName)
	}

	if mUser.Email != mUser.Email {
		uf.AddAttributeField("email", mUser.Email)
	}

	if mUser.Password != mUser.Password {
		uf.AddAttributeField("password", mUser.Password)
	}

	if mUser.PasswordSalt != mUser.PasswordSalt {
		uf.AddAttributeField("password_salt", mUser.PasswordSalt)
	}

	if mUser.PasswordHash != mUser.PasswordHash {
		uf.AddAttributeField("password_hash", mUser.PasswordHash)
	}

	if mUser.StateId != mUser.StateId {
		uf.AddAttributeField("state_id", mUser.StateId)
	}

	if mUser.UpdatedAt != mUser.UpdatedAt {
		uf.AddAttributeField("updated_at", mUser.UpdatedAt)
	}

	//Add keys
	uf.AddKeyField("id", key.Id)

	fs, fargs, ks, kargs := uf.FieldNameAndArgs()

	query = sqlhelpers.BuildUpdateQuery(uf.Table(), fs, ks)
	args = append(fargs, kargs...)

	return

}
