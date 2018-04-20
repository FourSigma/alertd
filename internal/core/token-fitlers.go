package core

type TokenFilter interface {
	OK(*Token) bool
	Valid() error
}
