package http

import "github.com/go-chi/chi"

type httpCtxKey string

const (
	CtxUserId httpCtxKey = "UserId"
)

var rootRoute = chi.NewRouter()
