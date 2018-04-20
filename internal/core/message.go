package core

import (
	"context"
	"time"

	"github.com/FourSigma/alertd/pkg/util"
	uuid "github.com/satori/go.uuid"
)

type MessageRepo interface {
	Create(context.Context, *Message) error
	Delete(context.Context, MessageKey) error
	List(context.Context, MessageFilter, ...Opts) (MessageList, error)
	Get(context.Context, MessageKey) (*Message, error)
	Update(context.Context, MessageKey, *Message) error
}

type MessageTypeId string
type Message struct {
	Id        uuid.UUID
	TopicId   uuid.UUID
	Msg       string
	TypeId    MessageTypeId
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u Message) Key() MessageKey {
	return MessageKey{
		Id: u.Id,
	}
}

func (u Message) TopicKey() TopicKey {
	return TopicKey{
		Id: u.TopicId,
	}
}

type MessageKey struct {
	Id uuid.UUID
}

func (u MessageKey) FieldSet() util.FieldSet {
	return util.NewFieldSet("MessageKey",
		util.NewField("Id", u.Id, &u.Id, false),
	)
}

type MessageList []*Message

func (u MessageList) Map() (m map[MessageKey]*Message) {
	m = map[MessageKey]*Message{}
	for _, v := range u {
		m[v.Key()] = v
	}
	return
}

// func (u MessageList) Reslove(tl List) {
// 	m := tl.Map()
// 	for _, v := range tl {
// 		if usr, ok := m[v.]; ok {
// 			usr.List = append(usr.List, v)
// 		}
// 	}

// 	return
// }

func (u MessageList) Update(fn func(*Message)) (rs MessageList) {
	for _, v := range u {
		fn(v)
	}
	return
}

func (u MessageList) Filter(filt MessageFilter) (rs MessageList) {
	rs = make([]*Message, len(u))
	for _, v := range u {
		if filt.OK(v) {
			rs = append(rs, v)
		}
	}
	return
}

func (u MessageList) KeyList() (kl []MessageKey) {
	kl = make([]MessageKey, len(u))
	for i, v := range u {
		kl[i] = v.Key()
	}
	return
}
