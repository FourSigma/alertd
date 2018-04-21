package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/FourSigma/alertd/internal/core"
	"github.com/FourSigma/alertd/internal/service"
	utilhttp "github.com/FourSigma/alertd/pkg/util/http"
	"github.com/go-chi/chi"
)

func TokenCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		usrKey := r.Context().Value(CtxUserId).(core.UserKey)

		tokenId := chi.URLParam(r, "tokenId")
		if tokenId == "" {
			utilhttp.HandleError(w, utilhttp.ErrorDecodingPathTokenId, nil)
			return
		}
		ctx := context.WithValue(r.Context(), CtxTokenId, core.TokenKey{Token: tokenId, UserId: usrKey.Id})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type TokenResource struct {
	token *service.TokenService
}

func (u TokenResource) Create(rw http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		utilhttp.HandleError(rw, utilhttp.ErrorEmptyBody, nil)
		return
	}
	usrKey := r.Context().Value(CtxUserId).(core.UserKey)

	req := &service.TokenCreateRequest{Data: &core.Token{UserId: usrKey.Id}}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		utilhttp.HandleError(rw, utilhttp.ErrorJSONDecoding, err)
		return
	}
	resp, err := u.token.Create(r.Context(), req)
	if err != nil {
		utilhttp.HandleError(rw, utilhttp.ErrorCreatingResource, err)
		return
	}
	utilhttp.JSONResponse(rw, http.StatusCreated, &utilhttp.Response{Data: resp.Data})
}

func (u TokenResource) Get(rw http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(CtxTokenId).(core.TokenKey)
	resp, err := u.token.Get(r.Context(), &service.TokenGetRequest{Key: key})
	if err != nil {
		utilhttp.HandleError(rw, utilhttp.ErrorGetResource, err)
		return
	}
	utilhttp.JSONResponse(rw, http.StatusOK, &utilhttp.Response{Data: resp.Data})
}

func (u TokenResource) Delete(rw http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(CtxTokenId).(core.TokenKey)
	resp, err := u.token.Delete(r.Context(), &service.TokenDeleteRequest{Key: key})
	if err != nil {
		utilhttp.HandleError(rw, utilhttp.ErrorDeletingResource, err)
		return
	}
	utilhttp.JSONResponse(rw, http.StatusOK, &utilhttp.Response{Data: resp.Key})
}

func (u TokenResource) Update(rw http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		utilhttp.HandleError(rw, utilhttp.ErrorEmptyBody, nil)
		return
	}
	key := r.Context().Value(CtxTokenId).(core.TokenKey)
	req := &service.TokenUpdateRequest{Key: key}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		utilhttp.HandleError(rw, utilhttp.ErrorJSONDecoding, err)
		return
	}

	resp, err := u.token.Update(r.Context(), req)
	if err != nil {
		utilhttp.HandleError(rw, utilhttp.ErrorUpdatingResource, err)
		return
	}
	utilhttp.JSONResponse(rw, http.StatusOK, &utilhttp.Response{Data: resp.Data})
}

func (u TokenResource) Index(rw http.ResponseWriter, r *http.Request) {
}
