package service

import (
	"context"
	"fmt"

	"github.com/malkev1ch/observability/services/apiservice/internal/model"
)

type VoucherRepository interface {
	GetByUserID(ctx context.Context, userID int64) (*model.Voucher, error)
}

type VoucherUserRepository interface {
	GetByID(ctx context.Context, userID int64) (*model.User, error)
}

type Voucher struct {
	rps  VoucherRepository
	uRps VoucherUserRepository
}

func NewVoucher(rps VoucherRepository, uRps VoucherUserRepository) *Voucher {
	return &Voucher{
		rps:  rps,
		uRps: uRps,
	}
}

func (s *Voucher) Search(ctx context.Context, userID int64) (*model.Voucher, error) {
	voucher, err := s.rps.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get voucher by user id: %w", err)
	}

	user, err := s.uRps.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	voucher.UserName = user.Name

	return voucher, nil
}
