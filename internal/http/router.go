package http

import (
	"net/http"
	"time"

	"github.com/FourSigma/alertd/internal/repo/postgres"
	"github.com/FourSigma/alertd/internal/service"
	utilhttp "github.com/FourSigma/alertd/pkg/util/http"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	log "github.com/Sirupsen/logrus"
)

type httpCtxKey string

const (
	CtxUserId    httpCtxKey = "UserKey"
	CtxTokenId   httpCtxKey = "TokenKey"
	CtxTopicId   httpCtxKey = "TopicKey"
	CtxMessageId httpCtxKey = "MessageKey"
)

func RepoCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := postgres.AddRepoContext(r.Context())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NewAPI(port string, l *log.Logger) *api {
	return &api{
		User: &UserResource{
			user: service.NewUserService(l),
		},
		Token: &TokenResource{
			token: service.NewTokenService(l),
		},
		Topic: &TopicResource{
			topic: service.NewTopicService(l),
		},
		Message: &MessageResource{
			message: service.NewMessageService(l),
		},
		log: l.WithFields(log.Fields{
			"layer": "api",
		}),
		r:    chi.NewRouter(),
		port: port,
	}
}

type api struct {
	User    *UserResource
	Token   *TokenResource
	Topic   *TopicResource
	Message *MessageResource

	log *log.Entry

	r    chi.Router
	port string
}

func (u *api) Run() error {
	return http.ListenAndServe(":"+u.port, u.routes())
}

func RequestLogger(l *log.Entry) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			defer func() {
				entry := l.WithFields(log.Fields{
					"method": r.Method,
					"status": ww.Status(),
					"bytes":  ww.BytesWritten(),
					"time":   time.Since(t1),
					"reqId":  middleware.GetReqID(r.Context()),
				})
				if ww.Status() >= 300 && ww.Status() < 500 {
					entry.Warnf("%s %s", r.Method, r.URL.String())
					return
				}

				if ww.Status() >= 500 {
					entry.Errorf("%s %s", r.Method, r.URL.String())
					return
				}
				entry.Infof("%s %s", r.Method, r.URL.String())
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}

}

func (u *api) routes() chi.Router {
	//Middleware
	u.r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Recoverer,
		RepoCtx,
		RequestLogger(u.log),
	)
	u.r.Route("/v1", func(r chi.Router) {
		u.userRoutes(r)
	})
	return u.r
}

func (u *api) userRoutes(r chi.Router) {
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

			r.Route("/topics", func(r chi.Router) {
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

func (u *api) tokenRoutes(r chi.Router) {
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

func (u *api) topicRoutes(r chi.Router) {
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
