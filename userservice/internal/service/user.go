package service

import (
	"context"
	"fmt"

	"github.com/malkev1ch/observability/userservice/internal/model"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
	Create(ctx context.Context, user *model.User) (*model.User, error)
}

type VoucherRepository interface {
	Create(ctx context.Context, userID int64) error
}

type User struct {
	rps  UserRepository
	vRps VoucherRepository
}

func (s *User) GetByID(ctx context.Context, id int64) (*model.User, error) {
	return s.rps.GetByID(ctx, id)
}

func (s *User) Create(ctx context.Context, user *model.User) (*model.User, error) {
	user, err := s.rps.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	err = s.vRps.Create(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create user's voucher: %w", err)
	}

	return user, nil
}

func NewUser(rps UserRepository, vRps VoucherRepository) *User {
	return &User{rps: rps, vRps: vRps}
}
