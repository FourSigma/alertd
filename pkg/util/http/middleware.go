package http

import (
	"net/http"
)

func ParseQuery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			HandleError(w, ErrorParsingQueryParams, err)
			return
		}
		//ctx := context.WithValue(r.Context(), CtxUserId, core.UserKey{Id: userId})
		//next.ServeHTTP(w, r.WithContext(ctx))
	})
}
