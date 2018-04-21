package http

import (
	"github.com/go-chi/chi"
)

type httpCtxKey string

const (
	CtxUserId    httpCtxKey = "UserId"
	CtxTokenId   httpCtxKey = "TokenId"
	CtxTopicId   httpCtxKey = "TopicId"
	CtxMessageId httpCtxKey = "MessageId"
)

var rootRoute = chi.NewRouter()

type ResourceManager struct {
	User  *UserResource
	Token *TokenResource
}

var res *ResourceManager

func GetResourceManager() *ResourceManager {
	if res == nil {
		res = NewResourceManager()
		return res
	}
	return res
}

func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		User:  &UserResource{},
		Token: &TokenResource{},
	}
}
