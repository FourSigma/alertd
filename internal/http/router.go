package http

import (
	"net/http"

	"github.com/FourSigma/alertd/internal/repo/postgres"
	"github.com/FourSigma/alertd/internal/service"
	utilhttp "github.com/FourSigma/alertd/pkg/util/http"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type httpCtxKey string

const (
	CtxUserId    httpCtxKey = "UserId"
	CtxTokenId   httpCtxKey = "TokenId"
	CtxTopicId   httpCtxKey = "TopicId"
	CtxMessageId httpCtxKey = "MessageId"
)

func RepoCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := postgres.AddRepoContext(r.Context())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

var rootRoute = chi.NewRouter().With(RepoCtx)

func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		User: &UserResource{
			user: service.NewUserService(),
		},
		Token: &TokenResource{
			token: service.NewTokenService(),
		},
		Topic: &TopicResource{
			topic: service.NewTopicService(),
		},
		Message: &MessageResource{
			message: service.NewMessageService(),
		},
	}
}

type ResourceManager struct {
	User    *UserResource
	Token   *TokenResource
	Topic   *TopicResource
	Message *MessageResource
}

func (u *ResourceManager) Routes() (r chi.Router) {
	r = chi.NewRouter()
	//Middleware
	r.Use(
		middleware.Logger,

		RepoCtx,
	)
	return r.Route("/v1", func(r chi.Router) {
		u.userRoutes(r)
	})
}

func (u *ResourceManager) userRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", u.User.Create)
		r.With(utilhttp.ParseQuery).Get("/", u.User.Index)

		r.Route("/{userId}", func(r chi.Router) {
			r.Use(UserCtx)
			r.Get("/", u.User.Get)
			r.Put("/", u.User.Update)
			r.Delete("/", u.User.Delete)

			r.Route("/tokens", func(r chi.Router) {
				r.Post("/", u.Token.Create)
				r.Route("/{tokenId}", func(r chi.Router) {
					r.Use(TokenCtx)
					r.Get("/", u.Token.Get)
					r.Put("/", u.Token.Update)
					r.Delete("/", u.Token.Delete)
				})
			})
		})
	})
}
