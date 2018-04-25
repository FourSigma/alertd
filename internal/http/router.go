package http

import (
	"context"
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
	CtxUserEntity httpCtxKey = "core.User"
	CtxUserId     httpCtxKey = "UserKey"
	CtxTokenId    httpCtxKey = "TokenKey"
	CtxTopicId    httpCtxKey = "TopicKey"
	CtxMessageId  httpCtxKey = "MessageKey"
)

func RepoCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := postgres.AddRepoContext(r.Context())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NewAPI(port string, l *log.Logger) *api {
	svc := service.GetService(l)
	return &api{
		User: UserResource{
			user: svc.User,
		},
		Token: TokenResource{
			token: svc.Token,
		},
		Topic: TopicResource{
			topic: svc.Topic,
		},
		Message: MessageResource{
			message: svc.Message,
		},
		log: l.WithFields(log.Fields{
			"layer": "api",
		}),
		r:    chi.NewRouter(),
		port: port,
	}
}

type api struct {
	User    UserResource
	Token   TokenResource
	Topic   TopicResource
	Message MessageResource

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

func (u *api) TokenCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Token")
		if token == "" {
			utilhttp.HandleError(w, utilhttp.ErrorTokenRequired, nil)
			return
		}
		usr, err := u.User.user.GetUserFromToken(r.Context(), token)
		if err != nil {
			utilhttp.HandleError(w, utilhttp.ErrorDecodingPathTokenId, nil)
			return
		}
		ctx := context.WithValue(r.Context(), CtxUserEntity, usr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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

			// r.Route("/tokens", func(r chi.Router) {
			// 	r.Post("/", u.Token.Create)
			// 	r.Route("/{tokenId}", func(r chi.Router) {
			// 		r.Use(TokenCtx)
			// 		r.Get("/", u.Token.Get)
			// 		r.Put("/", u.Token.Update)
			// 		r.Delete("/", u.Token.Delete)
			// 	})
			// })

		})
	})
	r.Route("/topics", func(r chi.Router) {
		r.Use(u.TokenCtx)
		r.Get("/", u.Topic.Index)
		r.Post("/", u.Topic.Create)
		r.Route("/{topicId}", func(r chi.Router) {
			r.Use(TopicCtx)
			r.Get("/", u.Topic.Get)
			r.Put("/", u.Topic.Update)
			r.Delete("/", u.Topic.Delete)
		})
	})
}

func (u *api) tokenRoutes(r chi.Router) {
	// r.Route("/tokens", func(r chi.Router) {
	// 	r.With(utilhttp.ParseQuery).Get("/", u.Token.Index)

	// 	r.Route("/{tokenId}", func(r chi.Router) {
	// 		r.Use(UserCtx)
	// 		r.Get("/", u.User.Get)
	// 		r.Put("/", u.User.Update)
	// 		r.Delete("/", u.User.Delete)

	// 		r.Route("/tokens", func(r chi.Router) {
	// 			r.Post("/", u.Token.Create)
	// 			r.Route("/{tokenId}", func(r chi.Router) {
	// 				r.Use(TokenCtx)
	// 				r.Get("/", u.Token.Get)
	// 				r.Put("/", u.Token.Update)
	// 				r.Delete("/", u.Token.Delete)
	// 			})
	// 		})
	// 	})
	// })
}

func (u *api) topicRoutes(r chi.Router) {
	// r.Route("/tokens", func(r chi.Router) {
	// 	r.With(utilhttp.ParseQuery).Get("/", u.Token.Index)

	// 	r.Route("/{tokenId}", func(r chi.Router) {
	// 		r.Use(UserCtx)
	// 		r.Get("/", u.User.Get)
	// 		r.Put("/", u.User.Update)
	// 		r.Delete("/", u.User.Delete)

	// 		r.Route("/tokens", func(r chi.Router) {
	// 			r.Post("/", u.Token.Create)
	// 			r.Route("/{tokenId}", func(r chi.Router) {
	// 				r.Use(TokenCtx)
	// 				r.Get("/", u.Token.Get)
	// 				r.Put("/", u.Token.Update)
	// 				r.Delete("/", u.Token.Delete)
	// 			})
	// 		})
	// 	})
	// })
}
