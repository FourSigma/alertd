package service

import (
	"context"

	"github.com/FourSigma/alertd/internal/core"
	"github.com/FourSigma/alertd/internal/repo"
)

type TopicCreateRequest struct {
	Data core.Topic
}

type TopicCreateResponse struct {
	Data core.Topic
}
type TopicGetRequest struct {
	Key core.TopicKey
}
type TopicGetResponse struct {
	Data core.Topic
}

type TopicDeleteRequest struct {
	Key core.TopicKey
}
type TopicDeleteResponse struct {
	Key core.TopicKey
}

type TopicUpdateRequest struct {
	Key  core.TopicKey
	Data core.Topic
}
type TopicUpdateResponse struct {
	Data core.Topic
}

type TopicListRequest struct {
	Filter core.TopicFilter
	Opts   []core.Opts
}

type TopicListResponse struct {
	Data core.TopicList
}

type TopicService struct {
	repo repo.Datastore
}

func (u TopicService) Create(ctx context.Context, req *TopicCreateRequest) (resp *TopicCreateResponse, err error) {
	if err = u.repo.Topic.Create(ctx, &req.Data); err != nil {
		return
	}
	resp = &TopicCreateResponse{Data: req.Data}
	return
}
func (u TopicService) Update(ctx context.Context, req *TopicUpdateRequest) (resp *TopicUpdateResponse, err error) {
	if err = u.repo.Topic.Update(ctx, req.Key, &req.Data); err != nil {
		return
	}
	resp = &TopicUpdateResponse{Data: req.Data}
	return
}

func (u TopicService) Delete(ctx context.Context, req *TopicDeleteRequest) (resp *TopicDeleteResponse, err error) {
	if err = u.repo.Topic.Delete(ctx, req.Key); err != nil {
		return
	}
	resp = &TopicDeleteResponse{Key: req.Key}
	return
}

func (u TopicService) Get(ctx context.Context, req *TopicGetRequest) (resp *TopicGetResponse, err error) {
	resp = &TopicGetResponse{}
	if resp.Data, err = u.repo.Topic.Get(ctx, req.Key); err != nil {
		return
	}
	return
}

func (u TopicService) List(ctx context.Context, req *TopicListRequest) (resp *TopicListResponse, err error) {
	var ds []*core.Topic
	if ds, err = u.repo.Topic.List(ctx, req.Filter, req.Opts...); err != nil {
		return
	}
	resp = &TopicListResponse{Data: ds}
	return
}
