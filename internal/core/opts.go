package core

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
