package service

import (
	"context"
	"github.com/malkev1ch/observability/userservice/internal/model"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
	Create(ctx context.Context, user *model.User) (*model.User, error)
}

type User struct {
	rps UserRepository
}

func (s *User) GetByID(ctx context.Context, id int64) (*model.User, error) {
	return s.rps.GetByID(ctx, id)
}

func (s *User) Create(ctx context.Context, user *model.User) (*model.User, error) {
	return s.rps.Create(ctx, user)
}

func NewUser(rps UserRepository) *User {
	return &User{rps: rps}
}
