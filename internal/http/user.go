package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/FourSigma/alertd/internal/core"
	"github.com/FourSigma/alertd/internal/service"
	utilhttp "github.com/FourSigma/alertd/pkg/util/http"
	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"
)

func UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		strId := chi.URLParam(r, "userId")
		userId, err := uuid.FromString(strId)
		if err != nil {
			utilhttp.HandleError(w, utilhttp.ErrorDecodingPathUserId, err)
			return
		}
		ctx := context.WithValue(r.Context(), CtxUserId, core.UserKey{Id: userId})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type UserResource struct {
	user *service.UserService
}

func (u UserResource) Create(rw http.ResponseWriter, r *http.Request) {
	req := &service.UserCreateRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		utilhttp.HandleError(rw, utilhttp.ErrorJSONDecoding, err)
		return
	}
	resp, err := u.user.Create(r.Context(), req)
	if err != nil {
		utilhttp.HandleError(rw, utilhttp.ErrorCreatingResource, err)
		return
	}
	utilhttp.JSONResponse(rw, http.StatusCreated, &utilhttp.Response{Data: resp.Data})
}

func (u UserResource) Get(rw http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(CtxUserId).(core.UserKey)
	resp, err := u.user.Get(r.Context(), &service.UserGetRequest{Key: key})
	if err != nil {
		utilhttp.HandleError(rw, utilhttp.ErrorGetResource, err)
		return
	}
	utilhttp.JSONResponse(rw, http.StatusOK, &utilhttp.Response{Data: resp.Data})
}

func (u UserResource) Delete(rw http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(CtxUserId).(core.UserKey)
	resp, err := u.user.Delete(r.Context(), &service.UserDeleteRequest{Key: key})
	if err != nil {
		utilhttp.HandleError(rw, utilhttp.ErrorDeletingResource, err)
		return
	}
	utilhttp.JSONResponse(rw, http.StatusOK, &utilhttp.Response{Data: resp.Key})
}

func (u UserResource) Update(rw http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(CtxUserId).(core.UserKey)
	req := &service.UserUpdateRequest{Key: key}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		utilhttp.HandleError(rw, utilhttp.ErrorJSONDecoding, err)
		return
	}

	resp, err := u.user.Update(r.Context(), req)
	if err != nil {
		utilhttp.HandleError(rw, utilhttp.ErrorUpdatingResource, err)
		return
	}
	utilhttp.JSONResponse(rw, http.StatusOK, &utilhttp.Response{Data: resp.Data})
}

func (u UserResource) Index(rw http.ResponseWriter, r *http.Request) {
}
