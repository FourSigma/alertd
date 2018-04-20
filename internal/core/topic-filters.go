package core

type TopicFilter interface {
	OK(*Topic) bool
	Valid() error
}
