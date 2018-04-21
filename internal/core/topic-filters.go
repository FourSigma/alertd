package core

import "errors"

type TopicFilter interface {
	OK(*Topic) bool
	Valid() error
}

type FilterTopicAll struct{}

func (_ FilterTopicAll) OK(u *Topic) bool { return true }
func (_ FilterTopicAll) Valid() error     { return nil }

type FilterTopicKeyIn struct {
	KeyList []TopicKey

	cache map[TopicKey]struct{}
}

func (f *FilterTopicKeyIn) Valid() (err error) {
	if len(f.KeyList) == 0 {
		err = errors.New("Invalid FilterTopicKeyIn -- KeyList is empty")
		return
	}
	return nil
}

func (f *FilterTopicKeyIn) OK(u *Topic) (ok bool) {
	if f.cache != nil {
		f.cache = map[TopicKey]struct{}{}
		for _, v := range f.KeyList {
			f.cache[v] = struct{}{}
		}
	}
	_, ok = f.cache[u.Key()]
	return
}

type FilterTopicUserKeyIn struct {
	KeyList []UserKey

	cache map[UserKey]struct{}
}

func (f *FilterTopicUserKeyIn) Valid() (err error) {
	if len(f.KeyList) == 0 {
		err = errors.New("Invalid FilterTopicUserKeyIn -- KeyList is empty")
		return
	}
	return nil
}

func (f *FilterTopicUserKeyIn) OK(u *Topic) (ok bool) {
	if f.cache != nil {
		f.cache = map[UserKey]struct{}{}
		for _, v := range f.KeyList {
			f.cache[v] = struct{}{}
		}
	}
	_, ok = f.cache[u.UserKey()]
	return
}
