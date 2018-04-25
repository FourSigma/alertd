package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/FourSigma/alertd/internal/core"
	"github.com/FourSigma/alertd/internal/repo"
	"github.com/FourSigma/alertd/pkg/util"
)

type UserCreateRequest struct {
	core.User
	Password string `json:"password"`
}

type UserCreateResponse struct {
	Data core.User
}
type UserGetRequest struct {
	Key core.UserKey
}
type UserGetResponse struct {
	Data core.User
}

type UserDeleteRequest struct {
	Key core.UserKey
}
type UserDeleteResponse struct {
	Key core.UserKey
}

type UserUpdateRequest struct {
	Key core.UserKey
	core.User
}
type UserUpdateResponse struct {
	Data core.User
}

type UserListRequest struct {
	Filter core.UserFilter
	Opts   []core.Opts
}

type UserListResponse struct {
	Data core.UserList
}

type UserService struct {
	repo repo.Datastore
}

func (u UserService) Create(ctx context.Context, req *UserCreateRequest) (resp *UserCreateResponse, err error) {
	req.PasswordSalt, req.PasswordHash, err = util.EncryptPassword(req.Password)
	if err != nil {
		return
	}
	if err = u.repo.User.Create(ctx, &req.User); err != nil {
		fmt.Println(err)
		return
	}
	resp = &UserCreateResponse{Data: req.User}
	return
}
func (u UserService) Update(ctx context.Context, req *UserUpdateRequest) (resp *UserUpdateResponse, err error) {
	if err = u.repo.User.Update(ctx, req.Key, &req.User); err != nil {
		fmt.Println(err)
		return
	}
	resp = &UserUpdateResponse{Data: req.User}
	return
}

func (u UserService) Delete(ctx context.Context, req *UserDeleteRequest) (resp *UserDeleteResponse, err error) {
	if err = u.repo.User.Delete(ctx, req.Key); err != nil {
		return
	}
	resp = &UserDeleteResponse{Key: req.Key}
	return
}

func (u UserService) Get(ctx context.Context, req *UserGetRequest) (resp *UserGetResponse, err error) {
	var d core.User
	if d, err = u.repo.User.Get(ctx, req.Key); err != nil {
		return
	}
	resp = &UserGetResponse{Data: d}
	return
}

func (u UserService) List(ctx context.Context, req *UserListRequest) (resp *UserListResponse, err error) {
	var ds []*core.User
	if ds, err = u.repo.User.List(ctx, req.Filter, req.Opts...); err != nil {
		return
	}
	resp = &UserListResponse{Data: ds}
	return
}
func (u UserService) GetUserFromToken(ctx context.Context, token string) (usr core.User, err error) {
	var ds core.TokenList
	if ds, err = u.repo.Token.List(ctx, &core.FilterTokenIn{TokenList: []string{token}}); err != nil {
		return
	}
	if len(ds) > 1 || len(ds) == 0 {
		err = errors.New("zero or multiple users for this token")
		return
	}

	if usr, err = u.repo.User.Get(ctx, ds[0].UserKey()); err != nil {
		return
	}
	return

}

func (u UserService) AuthList(ctx context.Context, req *UserListRequest) (err error) {
	err = ErrorUnauthorized{}
	return
}

func (u UserService) AuthGet(ctx context.Context, req *UserGetRequest) (err error) {
	err = ErrorUnauthorized{}
	return
}

func (u UserService) AuthDelete(ctx context.Context, req *UserDeleteRequest) (err error) {
	err = ErrorUnauthorized{}
	return
}

func (u UserService) AuthCreate(ctx context.Context, req *UserCreateRequest) (err error) {
	err = ErrorUnauthorized{}
	return
}

func (u UserService) AuthUpdate(ctx context.Context, req *UserUpdateRequest) (err error) {
	err = ErrorUnauthorized{}
	return
}

type ErrorUnauthorized struct{}

func (e ErrorUnauthorized) Error() string {
	return "Unauthorized Access"
}
