package handler

import (
	"context"
	v1 "github.com/malkev1ch/observability/userservice/gen/user/v1"
	"github.com/malkev1ch/observability/userservice/internal/model"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
	"time"
)

type UserService interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
	Create(ctx context.Context, user *model.User) (*model.User, error)
}

type User struct {
	svc UserService

	v1.UnimplementedUserServiceServer
}

func NewUser(svc UserService) *User {
	return &User{svc: svc}
}

// TODO: consider to inline in every place while building PGO

func ToUser(user *model.User) *v1.User {
	return &v1.User{
		Id:        user.ID,
		Name:      user.Name,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}

func (h *User) Create(ctx context.Context, in *v1.CreateRequest) (*v1.CreateResponse, error) {
	user, err := h.svc.Create(ctx, &model.User{
		ID:        0,
		Name:      in.Name,
		CreatedAt: time.Time{},
	})
	if err != nil {
		slog.Error(
			"Failed to create user",
			slog.Any("request", in),
		)
		return nil, err
	}

	return &v1.CreateResponse{User: ToUser(user)}, nil
}

func (h *User) GetByID(ctx context.Context, in *v1.GetByIDRequest) (*v1.GetByIDResponse, error) {
	user, err := h.svc.GetByID(ctx, in.Id)
	if err != nil {
		slog.Error(
			"Failed to get user",
			slog.Any("request", in),
		)
		return nil, err
	}

	return &v1.GetByIDResponse{User: ToUser(user)}, nil
}
