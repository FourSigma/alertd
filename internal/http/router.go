package http

import (
	"github.com/FourSigma/alertd/internal/service"
	"github.com/go-chi/chi"
)

type httpCtxKey string

const (
	CtxUserId  httpCtxKey = "UserId"
	CtxTokenId httpCtxKey = "TokenId"
)

var rootRoute = chi.NewRouter()

type ServiceManger struct {
	User service.UserService
}
