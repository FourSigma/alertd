package service

import (
	"context"

	"github.com/FourSigma/alertd/internal/core"
)

type UserCreateRequest struct {
	Data *core.User
}

type UserCreateResponse struct {
	Data *core.User
}
type UserGetRequest struct {
	Key core.UserKey
}
type UserGetResponse struct {
	Data *core.User
}

type UserDeleteRequest struct {
	Key core.UserKey
}
type UserDeleteResponse struct {
	Key core.UserKey
}

type UserUpdateRequest struct {
	Key  core.UserKey
	Data *core.User
}
type UserUpdateResponse struct {
	Data *core.User
}

type UserListRequest struct {
	Filter core.UserFilter
	Opts   []core.Opts
}

type UserListResponse struct {
	Data core.UserList
}

type UserService struct {
	usrRepo core.UserRepo
}

func (u UserService) Create(ctx context.Context, req *UserCreateRequest) (resp *UserCreateResponse, err error) {
	if err = u.usrRepo.Create(ctx, req.Data); err != nil {
		return
	}
	resp = &UserCreateResponse{Data: req.Data}
	return
}
func (u UserService) Update(ctx context.Context, req *UserUpdateRequest) (resp *UserUpdateResponse, err error) {
	if err = u.usrRepo.Update(ctx, req.Key, req.Data); err != nil {
		return
	}
	resp = &UserUpdateResponse{Data: req.Data}
	return
}

func (u UserService) Delete(ctx context.Context, req *UserDeleteRequest) (resp *UserDeleteResponse, err error) {
	if err = u.usrRepo.Delete(ctx, req.Key); err != nil {
		return
	}
	resp = &UserDeleteResponse{Key: req.Key}
	return
}

func (u UserService) Get(ctx context.Context, req *UserGetRequest) (resp *UserGetResponse, err error) {
	var d *core.User
	if d, err = u.usrRepo.Get(ctx, req.Key); err != nil {
		return
	}
	resp = &UserGetResponse{Data: d}
	return
}

func (u UserService) List(ctx context.Context, req *UserListRequest) (resp *UserListResponse, err error) {
	var ds []*core.User
	if ds, err = u.usrRepo.List(ctx, req.Filter, req.Opts...); err != nil {
		return
	}
	resp = &UserListResponse{Data: ds}
	return
}
