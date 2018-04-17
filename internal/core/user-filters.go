package core

type FilterUserAll struct{}

func (_ FilterUserAll) OK(u *User) bool { return true }

type FilterUserActiveUsers struct{}

func (_ FilterUserActiveUsers) OK(u *User) bool { return u.StateId == "Active" }

type FilterUserKeyIn struct {
	KeyList []UserKey

	cache map[UserKey]struct{}
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
