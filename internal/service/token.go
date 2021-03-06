package service

import (
	"context"

	"github.com/FourSigma/alertd/internal/core"
	"github.com/FourSigma/alertd/internal/repo"
)

type TokenCreateRequest struct {
	Data *core.Token
}

type TokenCreateResponse struct {
	Data *core.Token
}
type TokenGetRequest struct {
	Key core.TokenKey
}
type TokenGetResponse struct {
	Data core.Token
}

type TokenDeleteRequest struct {
	Key core.TokenKey
}
type TokenDeleteResponse struct {
	Key core.TokenKey
}

type TokenUpdateRequest struct {
	Key  core.TokenKey
	Data *core.Token
}
type TokenUpdateResponse struct {
	Data *core.Token
}

type TokenListRequest struct {
	Filter core.TokenFilter
	Opts   []core.Opts
}

type TokenListResponse struct {
	Data core.TokenList
}

type TokenService struct {
	repo repo.Datastore
}

func (u TokenService) Create(ctx context.Context, req *TokenCreateRequest) (resp *TokenCreateResponse, err error) {
	if err = u.repo.Token.Create(ctx, req.Data); err != nil {
		return
	}
	resp = &TokenCreateResponse{Data: req.Data}
	return
}
func (u TokenService) Update(ctx context.Context, req *TokenUpdateRequest) (resp *TokenUpdateResponse, err error) {
	if err = u.repo.Token.Update(ctx, req.Key, req.Data); err != nil {
		return
	}
	resp = &TokenUpdateResponse{Data: req.Data}
	return
}

func (u TokenService) Delete(ctx context.Context, req *TokenDeleteRequest) (resp *TokenDeleteResponse, err error) {
	if err = u.repo.Token.Delete(ctx, req.Key); err != nil {
		return
	}
	resp = &TokenDeleteResponse{Key: req.Key}
	return
}

func (u TokenService) Get(ctx context.Context, req *TokenGetRequest) (resp *TokenGetResponse, err error) {
	resp = &TokenGetResponse{}
	if resp.Data, err = u.repo.Token.Get(ctx, req.Key); err != nil {
		return
	}
	return
}

func (u TokenService) List(ctx context.Context, req *TokenListRequest) (resp *TokenListResponse, err error) {
	var ds []*core.Token
	if ds, err = u.repo.Token.List(ctx, req.Filter, req.Opts...); err != nil {
		return
	}
	resp = &TokenListResponse{Data: ds}
	return
}
