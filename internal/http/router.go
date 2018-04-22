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

func NewAPI() *API {
	return &API{
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
		r: chi.NewRouter(),
	}
}

type API struct {
	User    *UserResource
	Token   *TokenResource
	Topic   *TopicResource
	Message *MessageResource

	r chi.Router
}

func (u *API) Routes() (r chi.Router) {
	//Middleware
	u.r.Use(
		middleware.Logger,
		RepoCtx,
	)
	return u.r.Route("/v1", func(r chi.Router) {
		u.userRoutes(r)
	})
}

func (u *API) userRoutes(r chi.Router) {
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

func (u *API) tokenRoutes(r chi.Router) {
	r.Route("/tokens", func(r chi.Router) {
		r.With(utilhttp.ParseQuery).Get("/", u.Token.Index)

		r.Route("/{tokenId}", func(r chi.Router) {
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
