package client

import (
	"context"

	"github.com/malkev1ch/observability/services/apiservice/internal/model"
	userv1 "github.com/malkev1ch/observability/services/userservice/gen/user/v1"
)

type User struct {
	cln userv1.UserServiceClient
}

func NewUser(cln userv1.UserServiceClient) *User {
	return &User{cln: cln}
}

func FromUser(in *userv1.User) *model.User {
	return &model.User{
		ID:        in.Id,
		Name:      in.Name,
		CreatedAt: in.CreatedAt.AsTime(),
	}
}

func (r *User) GetByID(ctx context.Context, id int64) (*model.User, error) {
	resp, err := r.cln.GetByID(ctx, &userv1.GetByIDRequest{Id: id})
	if err != nil {
		return nil, err
	}

	return FromUser(resp.User), nil
}

func (r *User) Create(ctx context.Context, user *model.User) (*model.User, error) {
	resp, err := r.cln.Create(ctx, &userv1.CreateRequest{Name: user.Name})
	if err != nil {
		return nil, err
	}

	return FromUser(resp.User), nil
}
