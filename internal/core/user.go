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
	List(context.Context, UserFilter, ...Opts) ([]*User, error)
	Get(context.Context, UserKey) (*User, error)
	Update(context.Context, UserKey, *User) error
}

type UserFilter interface {
	UserFilter()
}

type Opts interface {
	Opts()
}

type Limit uint

func (_ Limit) Opts() {
	return
}

type Offset uint

func (_ Offset) Opts() {
	return
}

type UserList []*User

func (u UserList) KeyList() (kl []UserKey) {
	kl = make([]UserKey, len(u))
	for i, v := range u {
		kl[i] = v.Key()
	}
	return
}

type FilterUserActiveUsers struct{}

func (_ FilterUserActiveUsers) UserFilter() {}

type FilterUserKeyIn struct {
	KeyList []UserKey
}

func (_ FilterUserKeyIn) UserFilter() {}

type UserStateId string
type User struct {
	Id           uuid.UUID
	FirstName    string
	LastName     string
	Email        string
	PasswordSalt string
	PasswordHash string
	StateId      UserStateId

	CreatedAt time.Time
	UpdatedAt time.Time
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

type UserKey struct {
	Id uuid.UUID
}

func (f UserKey) Args() []interface{} {
	return []interface{}{
		f.Id,
	}
}

func (u *UserKey) FieldSet() util.FieldSet {
	return util.NewFieldSet("UserKey",
		util.NewField("Id", u.Id, &u.Id, false),
	)
}

//Users can only have one active key at a time.
//Tokens are unique
type UserTokenStateId string
type UserToken struct {
	UserId    uuid.UUID
	Token     string
	StateId   UserTokenStateId
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserTokenKey struct {
	UserId uuid.UUID
	Token  string
}

//Users can have many topics.
//Topic Key is UserId and Name
type Topic struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	Name      string //unique
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TopicKey struct {
	Id uuid.UUID
}

type MessageTypeId string
type TopicMessage struct {
	TopicId   uuid.UUID
	TypeId    MessageTypeId
	Msg       string
	CreatedAt time.Time
	UpdatedAt time.Time
}
