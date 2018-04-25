package core

import (
	"context"
	"time"

	"github.com/FourSigma/alertd/pkg/util"
	uuid "github.com/satori/go.uuid"
)

type TokenRepo interface {
	Create(context.Context, *Token) error
	Delete(context.Context, TokenKey) error
	List(context.Context, TokenFilter, ...Opts) (TokenList, error)
	Get(context.Context, TokenKey) (Token, error)
	Update(context.Context, TokenKey, *Token) error
}

//Users can only have one active key at a time.
//Tokens are unique
type TokenStateId string
type Token struct {
	UserId    uuid.UUID    `json:"userId,omitempty" db:"user_id"`
	Token     string       `json:"token,omitempty" db:"token"`
	StateId   TokenStateId `json:"stateId,omitempty" db:"state_id"`
	CreatedAt time.Time    `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt time.Time    `json:"updatedAt,omitempty" db:"updated_at"`

	User *User
}

func (u *Token) New() util.Entity {
	return &Token{}
}

func (u *Token) FieldSet() util.FieldSet {
	return util.NewFieldSet("Token",
		util.NewField("UserId", u.UserId, &u.UserId, false),
		util.NewField("Token", u.Token, &u.Token, true),
		util.NewField("StateId", u.StateId, &u.StateId, true),
		util.NewField("CreatedAt", u.CreatedAt, &u.CreatedAt, false),
		util.NewField("UpdatedAt", u.UpdatedAt, &u.UpdatedAt, true),
	)
}

func (u Token) Key() TokenKey {
	return TokenKey{
		UserId: u.UserId,
		Token:  u.Token,
	}
}

func (u Token) UserKey() UserKey {
	return UserKey{
		Id: u.UserId,
	}
}

type TokenKey struct {
	UserId uuid.UUID
	Token  string
}

func (u TokenKey) IsValid() error {
	return nil
}

func (u TokenKey) FieldSet() util.FieldSet {
	return util.NewFieldSet("TokenKey",
		util.NewField("UserId", u.UserId, &u.UserId, false),
		util.NewField("Token", u.Token, &u.Token, true),
	)
}

type TokenList []*Token

func (u TokenList) Map() (m map[TokenKey]*Token) {
	m = map[TokenKey]*Token{}
	for _, v := range u {
		m[v.Key()] = v
	}
	return
}

func (u TokenList) Reslove(ul UserList) {
	m := ul.Map()
	for _, v := range u {
		if usr, ok := m[v.UserKey()]; ok {
			v.User = usr
		}
	}
	return
}

func (u TokenList) Update(fn func(*Token)) (rs TokenList) {
	for _, v := range u {
		fn(v)
	}
	return
}

func (u TokenList) Filter(filt TokenFilter) (rs TokenList) {
	rs = make([]*Token, len(u))
	for _, v := range u {
		if filt.OK(v) {
			rs = append(rs, v)
		}
	}
	return
}

func (u TokenList) KeyList() (kl []TokenKey) {
	kl = make([]TokenKey, len(u))
	for i, v := range u {
		kl[i] = v.Key()
	}
	return
}
