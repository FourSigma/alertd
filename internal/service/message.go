package service

import (
	"context"

	"github.com/FourSigma/alertd/internal/core"
	"github.com/FourSigma/alertd/internal/repo"
)

type MessageCreateRequest struct {
	Data *core.Message
}

type MessageCreateResponse struct {
	Data *core.Message
}
type MessageGetRequest struct {
	Key core.MessageKey
}
type MessageGetResponse struct {
	Data core.Message
}

type MessageDeleteRequest struct {
	Key core.MessageKey
}
type MessageDeleteResponse struct {
	Key core.MessageKey
}

type MessageUpdateRequest struct {
	Key  core.MessageKey
	Data *core.Message
}
type MessageUpdateResponse struct {
	Data *core.Message
}

type MessageListRequest struct {
	Filter core.MessageFilter
	Opts   []core.Opts
}

type MessageListResponse struct {
	Data core.MessageList
}

type MessageService struct {
	repo repo.Datastore
}

func (u MessageService) Create(ctx context.Context, req *MessageCreateRequest) (resp *MessageCreateResponse, err error) {
	if err = u.repo.Message.Create(ctx, req.Data); err != nil {
		return
	}
	resp = &MessageCreateResponse{Data: req.Data}
	return
}
func (u MessageService) Update(ctx context.Context, req *MessageUpdateRequest) (resp *MessageUpdateResponse, err error) {
	if err = u.repo.Message.Update(ctx, req.Key, req.Data); err != nil {
		return
	}
	resp = &MessageUpdateResponse{Data: req.Data}
	return
}

func (u MessageService) Delete(ctx context.Context, req *MessageDeleteRequest) (resp *MessageDeleteResponse, err error) {
	if err = u.repo.Message.Delete(ctx, req.Key); err != nil {
		return
	}
	resp = &MessageDeleteResponse{Key: req.Key}
	return
}

func (u MessageService) Get(ctx context.Context, req *MessageGetRequest) (resp *MessageGetResponse, err error) {
	resp = &MessageGetResponse{}
	if resp.Data, err = u.repo.Message.Get(ctx, req.Key); err != nil {
		return
	}
	return
}

func (u MessageService) List(ctx context.Context, req *MessageListRequest) (resp *MessageListResponse, err error) {
	var ds []*core.Message
	if ds, err = u.repo.Message.List(ctx, req.Filter, req.Opts...); err != nil {
		return
	}
	resp = &MessageListResponse{Data: ds}
	return
}
