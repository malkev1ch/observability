package client

import (
	"context"

	"github.com/malkev1ch/observability/apiservice/internal/model"
	v1 "github.com/malkev1ch/observability/userservice/gen/user/v1"
)

type User struct {
	cln v1.UserServiceClient
}

func NewUser(cln v1.UserServiceClient) *User {
	return &User{cln: cln}
}

func FromUser(in *v1.User) *model.User {
	return &model.User{
		ID:        in.Id,
		Name:      in.Name,
		CreatedAt: in.CreatedAt.AsTime(),
	}
}

func (r *User) GetByID(ctx context.Context, id int64) (*model.User, error) {
	resp, err := r.cln.GetByID(ctx, &v1.GetByIDRequest{Id: id})
	if err != nil {
		return nil, err
	}

	return FromUser(resp.User), nil
}

func (r *User) Create(ctx context.Context, user *model.User) (*model.User, error) {
	resp, err := r.cln.Create(ctx, &v1.CreateRequest{Name: user.Name})
	if err != nil {
		return nil, err
	}

	return FromUser(resp.User), nil
}
