package core

type MessageFilter interface {
	OK(*Message) bool
	Valid() error
}
