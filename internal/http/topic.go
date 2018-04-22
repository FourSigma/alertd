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

func TopicCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		strId := chi.URLParam(r, "topicId")
		topicId, err := uuid.FromString(strId)
		if err != nil {
			utilhttp.HandleError(w, utilhttp.ErrorDecodingPathTopicId, err)
			return
		}
		ctx := context.WithValue(r.Context(), CtxTopicId, core.TopicKey{Id: topicId})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type TopicResource struct {
	topic *service.TopicService
}

func (u TopicResource) Create(rw http.ResponseWriter, r *http.Request) {
	req := &service.TopicCreateRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		utilhttp.HandleError(rw, utilhttp.ErrorJSONDecoding, err)
		return
	}
	resp, err := u.topic.Create(r.Context(), req)
	if err != nil {
		utilhttp.HandleError(rw, utilhttp.ErrorCreatingResource, err)
		return
	}
	utilhttp.JSONResponse(rw, http.StatusCreated, &utilhttp.Response{Data: resp.Data})
}

func (u TopicResource) Get(rw http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(CtxTopicId).(core.TopicKey)
	resp, err := u.topic.Get(r.Context(), &service.TopicGetRequest{Key: key})
	if err != nil {
		utilhttp.HandleError(rw, utilhttp.ErrorGetResource, err)
		return
	}
	utilhttp.JSONResponse(rw, http.StatusOK, &utilhttp.Response{Data: resp.Data})
}

func (u TopicResource) Delete(rw http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(CtxTopicId).(core.TopicKey)
	resp, err := u.topic.Delete(r.Context(), &service.TopicDeleteRequest{Key: key})
	if err != nil {
		utilhttp.HandleError(rw, utilhttp.ErrorDeletingResource, err)
		return
	}
	utilhttp.JSONResponse(rw, http.StatusOK, &utilhttp.Response{Data: resp.Key})
}

func (u TopicResource) Update(rw http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(CtxTopicId).(core.TopicKey)
	req := &service.TopicUpdateRequest{Key: key}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		utilhttp.HandleError(rw, utilhttp.ErrorJSONDecoding, err)
		return
	}

	resp, err := u.topic.Update(r.Context(), req)
	if err != nil {
		utilhttp.HandleError(rw, utilhttp.ErrorUpdatingResource, err)
		return
	}
	utilhttp.JSONResponse(rw, http.StatusOK, &utilhttp.Response{Data: resp.Data})
}

func (u TopicResource) Index(rw http.ResponseWriter, r *http.Request) {
}
