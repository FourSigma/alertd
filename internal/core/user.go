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
	Get(context.Context, UserKey) (*User, error)
	Update(context.Context, UserKey, *User) error
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

func NewUser(firstName string, lastName string, email string, password string) *User {
	salt, hash := util.EncryptPassword(password)
	return &User{
		Id:           uuid.NewV4(),
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordSalt: salt,
		PasswordHash: hash,
	}
}

type UserStateId string
type User struct {
	Id           uuid.UUID
	FirstName    string      `db:"first_name"`
	LastName     string      `db:"last_name"`
	Email        string      `db:"email"`
	PasswordSalt string      `db:"password_salt"`
	PasswordHash string      `db:"password_hash"`
	StateId      UserStateId `db:"state_id"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
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

func (u UserKey) FieldSet() util.FieldSet {
	return util.NewFieldSet("UserKey",
		util.NewField("Id", u.Id, &u.Id, false),
	)
}

func (_ UserKey) KeyLen() int {
	return 1
}

type UserList []*User

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
