package core

import (
	"context"
	"time"

	uuid "github.com/satori/go.uuid"
)

type TopicRepo interface {
	Create(context.Context, *Topic) error
	Delete(context.Context, TopicKey) error
	List(context.Context, TopicFilter, ...Opts) (TopicList, error)
	Get(context.Context, TopicKey) (*Topic, error)
	Update(context.Context, TopicKey, *Topic) error
}

//Users can have many topics.
//Topic Key is UserId and Name
type Topic struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	Name      string //unique
	CreatedAt time.Time
	UpdatedAt time.Time

	User *User
}

func (t Topic) Key() TopicKey {
	return TopicKey{
		Id: t.Id,
	}
}

func (t Topic) UserKey() UserKey {
	return UserKey{
		Id: t.UserId,
	}
}

type TopicKey struct {
	Id uuid.UUID
}

type TopicList []*Topic

func (u TopicList) Map() (m map[TopicKey]*Topic) {
	m = map[TopicKey]*Topic{}
	for _, v := range u {
		m[v.Key()] = v
	}
	return
}

func (u TopicList) ResloveTopic() (m map[TopicKey]*Topic) {
	m = map[TopicKey]*Topic{}
	for _, v := range u {
		m[v.Key()] = v
	}
	return
}

func (u TopicList) Update(fn func(*Topic)) (rs TopicList) {
	for _, v := range u {
		fn(v)
	}
	return
}

func (u TopicList) Filter(filt TopicFilter) (rs TopicList) {
	rs = make([]*Topic, len(u))
	for _, v := range u {
		if filt.OK(v) {
			rs = append(rs, v)
		}
	}
	return
}

func (u TopicList) KeyList() (kl []TopicKey) {
	kl = make([]TopicKey, len(u))
	for i, v := range u {
		kl[i] = v.Key()
	}
	return
}
