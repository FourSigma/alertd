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

func MessageCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		strId := chi.URLParam(r, "messageId")
		messageId, err := uuid.FromString(strId)
		if err != nil {
			utilhttp.HandleError(w, utilhttp.ErrorDecodingPathMessageId, err)
			return
		}
		ctx := context.WithValue(r.Context(), CtxMessageId, core.MessageKey{Id: messageId})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type MessageResource struct {
	message service.MessageService
}

func (u MessageResource) Create(rw http.ResponseWriter, r *http.Request) {
	req := &service.MessageCreateRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		utilhttp.HandleError(rw, utilhttp.ErrorJSONDecoding, err)
		return
	}
	resp, err := u.message.Create(r.Context(), req)
	if err != nil {
		utilhttp.HandleError(rw, utilhttp.ErrorCreatingResource, err)
		return
	}
	utilhttp.JSONResponse(rw, http.StatusCreated, &utilhttp.Response{Data: resp.Data})
}

func (u MessageResource) Get(rw http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(CtxMessageId).(core.MessageKey)
	resp, err := u.message.Get(r.Context(), &service.MessageGetRequest{Key: key})
	if err != nil {
		utilhttp.HandleError(rw, utilhttp.ErrorGetResource, err)
		return
	}
	utilhttp.JSONResponse(rw, http.StatusOK, &utilhttp.Response{Data: resp.Data})
}

func (u MessageResource) Delete(rw http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(CtxMessageId).(core.MessageKey)
	resp, err := u.message.Delete(r.Context(), &service.MessageDeleteRequest{Key: key})
	if err != nil {
		utilhttp.HandleError(rw, utilhttp.ErrorDeletingResource, err)
		return
	}
	utilhttp.JSONResponse(rw, http.StatusOK, &utilhttp.Response{Data: resp.Key})
}

func (u MessageResource) Update(rw http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(CtxMessageId).(core.MessageKey)
	req := &service.MessageUpdateRequest{Key: key}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		utilhttp.HandleError(rw, utilhttp.ErrorJSONDecoding, err)
		return
	}

	resp, err := u.message.Update(r.Context(), req)
	if err != nil {
		utilhttp.HandleError(rw, utilhttp.ErrorUpdatingResource, err)
		return
	}
	utilhttp.JSONResponse(rw, http.StatusOK, &utilhttp.Response{Data: resp.Data})
}

func (u MessageResource) Index(rw http.ResponseWriter, r *http.Request) {
}
