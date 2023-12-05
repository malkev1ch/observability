package handler

import (
	"context"
	"log/slog"
	"time"

	userv1 "github.com/malkev1ch/observability/services/userservice/gen/user/v1"
	"github.com/malkev1ch/observability/services/userservice/internal/model"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
	Create(ctx context.Context, user *model.User) (*model.User, error)
}

type User struct {
	svc UserService

	userv1.UnimplementedUserServiceServer
}

func NewUser(svc UserService) *User {
	return &User{svc: svc}
}

func ToUser(user *model.User) *userv1.User {
	return &userv1.User{
		Id:        user.ID,
		Name:      user.Name,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}

func (h *User) Create(ctx context.Context, in *userv1.CreateRequest) (*userv1.CreateResponse, error) {
	user, err := h.svc.Create(ctx, &model.User{
		ID:        0,
		Name:      in.Name,
		CreatedAt: time.Time{},
	})
	if err != nil {
		slog.Error(
			"Failed to create voucher",
			slog.Any("request", in),
		)
		return nil, err
	}

	return &userv1.CreateResponse{User: ToUser(user)}, nil
}

func (h *User) GetByID(ctx context.Context, in *userv1.GetByIDRequest) (*userv1.GetByIDResponse, error) {
	user, err := h.svc.GetByID(ctx, in.Id)
	if err != nil {
		slog.Error(
			"Failed to get voucher by id",
			slog.Any("request", in),
		)
		return nil, err
	}

	return &userv1.GetByIDResponse{User: ToUser(user)}, nil
}
