package core

import (
	"context"
	"time"

	"github.com/FourSigma/alertd/pkg/util"
	uuid "github.com/satori/go.uuid"
)

type UserRepo interface {
	Create(context.Context, *User) error
	Delete(context.Context, UserKey) error
	List(context.Context, UserFilter, ...Opts) (UserList, error)
	Get(context.Context, UserKey) (User, error)
	Update(context.Context, UserKey, *User) error
}

type UserStateId string
type User struct {
	Id           uuid.UUID   `json:"id,omitempty" db:"id"`
	FirstName    string      `json:"firstName,omitempty" db:"first_name"`
	LastName     string      `json:"lastName,omitempty" db:"last_name"`
	Email        string      `json:"email,omitempty" db:"email"`
	PasswordSalt string      `json:"-" db:"password_salt"`
	PasswordHash string      `json:"-" db:"password_hash"`
	StateId      UserStateId `json:"stateId,omitempty" db:"state_id"`
	CreatedAt    time.Time   `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt    time.Time   `json:"updatedAt,omitempty" db:"updated_at"`

	TopicList TopicList `json:"topicList,omitempty"`
	TokenList TokenList `json:"tokenList,omitempty"`
}

func (u User) Key() UserKey {
	return UserKey{
		Id: u.Id,
	}
}

func (u *User) FieldSet() util.FieldSet {
	return util.NewFieldSet("User",
		util.NewField("Id", u.Id, &u.Id, false),
		util.NewField("FirstName", u.FirstName, &u.FirstName, true),
		util.NewField("LastName", u.LastName, &u.LastName, true),
		util.NewField("Email", u.Email, &u.Email, true),
		util.NewField("PasswordSalt", u.PasswordSalt, &u.PasswordSalt, true),
		util.NewField("PasswordHash", u.PasswordHash, &u.PasswordHash, true),
		util.NewField("StateId", u.StateId, &u.StateId, true),
		util.NewField("CreatedAt", u.CreatedAt, &u.CreatedAt, false),
		util.NewField("UpdatedAt", u.UpdatedAt, &u.UpdatedAt, true),
	)
}

func (u *User) New() util.Entity {
	return &User{}
}

type UserKey struct {
	Id uuid.UUID
}

func (u UserKey) IsValid() error {
	return nil
}

func (u UserKey) FieldSet() util.FieldSet {
	return util.NewFieldSet("UserKey",
		util.NewField("Id", u.Id, &u.Id, false),
	)
}

type UserList []*User

func (u UserList) Map() (m map[UserKey]*User) {
	m = map[UserKey]*User{}
	for _, v := range u {
		m[v.Key()] = v
	}
	return
}

func (u UserList) Reslove(tl TopicList, utl TokenList) {
	m := u.Map()
	for _, v := range tl {
		if usr, ok := m[v.UserKey()]; ok {
			usr.TopicList = append(usr.TopicList, v)
		}
	}

	for _, v := range utl {
		if usr, ok := m[v.UserKey()]; ok {
			usr.TokenList = append(usr.TokenList, v)
		}
	}
	return
}

func (u UserList) Update(fn func(*User)) (rs UserList) {
	for _, v := range u {
		fn(v)
	}
	return
}

func (u UserList) Filter(filt UserFilter) (rs UserList) {
	rs = make([]*User, len(u))
	for _, v := range u {
		if filt.OK(v) {
			rs = append(rs, v)
		}
	}
	return
}

func (u UserList) KeyList() (kl []UserKey) {
	kl = make([]UserKey, len(u))
	for i, v := range u {
		kl[i] = v.Key()
	}
	return
}
