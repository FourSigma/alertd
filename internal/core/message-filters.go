package core

import "errors"

type MessageFilter interface {
	OK(*Message) bool
	Valid() error
}

type FilterMessageAll struct{}

func (_ FilterMessageAll) OK(u *Message) bool { return true }
func (_ FilterMessageAll) Valid() error       { return nil }

type FilterMessageByTypeId struct {
	TypeId MessageTypeId
}

func (_ FilterMessageByTypeId) OK(u *Message) bool { return u.TypeId == u.TypeId }
func (_ FilterMessageByTypeId) Valid() error       { return nil }

type FilterMessageKeyIn struct {
	KeyList []MessageKey

	cache map[MessageKey]struct{}
}

func (f *FilterMessageKeyIn) Valid() (err error) {
	if len(f.KeyList) == 0 {
		err = errors.New("Invalid FilterMessageKeyIn -- KeyList is empty")
		return
	}
	return nil
}

func (f *FilterMessageKeyIn) OK(u *Message) (ok bool) {
	if f.cache == nil {
		f.cache = map[MessageKey]struct{}{}
		for _, v := range f.KeyList {
			f.cache[v] = struct{}{}
		}
	}
	_, ok = f.cache[u.Key()]
	return
}
