package core

import "errors"

type TokenFilter interface {
	OK(*Token) bool
	Valid() error
}

type FilterTokenAll struct{}

func (_ FilterTokenAll) OK(u *Token) bool { return true }
func (_ FilterTokenAll) Valid() error     { return nil }

type FilterTokenByStateId struct {
	StateId TokenStateId
}

func (_ FilterTokenByStateId) OK(u *Token) bool { return u.StateId == u.StateId }
func (_ FilterTokenByStateId) Valid() error     { return nil }

type FilterTokenKeyIn struct {
	KeyList []TokenKey

	cache map[TokenKey]struct{}
}

func (f *FilterTokenKeyIn) Valid() (err error) {
	if len(f.KeyList) == 0 {
		err = errors.New("Invalid FilterTokenKeyIn -- KeyList is empty")
		return
	}
	return nil
}

func (f *FilterTokenKeyIn) OK(u *Token) (ok bool) {
	if f.cache != nil {
		f.cache = map[TokenKey]struct{}{}
		for _, v := range f.KeyList {
			f.cache[v] = struct{}{}
		}
	}
	_, ok = f.cache[u.Key()]
	return
}

type FilterTokenUserKeyIn struct {
	KeyList []UserKey

	cache map[UserKey]struct{}
}

func (f *FilterTokenUserKeyIn) Valid() (err error) {
	if len(f.KeyList) == 0 {
		err = errors.New("Invalid FilterTokenUserKeyIn -- KeyList is empty")
		return
	}
	return nil
}

func (f *FilterTokenUserKeyIn) OK(u *Token) (ok bool) {
	if f.cache != nil {
		f.cache = map[UserKey]struct{}{}
		for _, v := range f.KeyList {
			f.cache[v] = struct{}{}
		}
	}
	_, ok = f.cache[u.UserKey()]
	return
}
