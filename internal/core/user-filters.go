package core

import "errors"

type UserFilter interface {
	OK(*User) bool
	Valid() error
}

type FilterUserAll struct{}

func (_ FilterUserAll) OK(u *User) bool { return true }
func (_ FilterUserAll) Valid() error    { return nil }

type FilterUserByStateId struct {
	StateId UserStateId
}

func (_ FilterUserByStateId) OK(u *User) bool { return u.StateId == u.StateId }
func (_ FilterUserByStateId) Valid() error    { return nil }

type FilterUserKeyIn struct {
	KeyList []UserKey

	cache map[UserKey]struct{}
}

func (f FilterUserKeyIn) Valid() (err error) {
	if len(f.KeyList) == 0 {
		err = errors.New("Invalid FilterUserKeyIn -- KeyList is empty")
		return
	}
	return nil
}

func (f *FilterUserKeyIn) OK(u *User) (ok bool) {
	if f.cache != nil {
		f.cache = map[UserKey]struct{}{}
		for _, v := range f.KeyList {
			f.cache[v] = struct{}{}
		}
	}
	_, ok = f.cache[u.Key()]
	return
}
