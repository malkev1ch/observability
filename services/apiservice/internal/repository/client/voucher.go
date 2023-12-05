package client

import (
	"context"
	"fmt"

	"github.com/malkev1ch/observability/services/apiservice/internal/model"
	userv1 "github.com/malkev1ch/observability/services/voucherservice/gen/voucher/v1"
)

type Voucher struct {
	cln userv1.VoucherServiceClient
}

func NewVoucher(cln userv1.VoucherServiceClient) *Voucher {
	return &Voucher{cln: cln}
}

func FromVoucher(in *userv1.Voucher) *model.Voucher {
	return &model.Voucher{
		ID:        in.Id,
		Value:     in.Value,
		UserID:    in.UserId,
		UserName:  "",
		CreatedAt: in.CreatedAt.AsTime(),
	}
}

func (r *Voucher) GetByUserID(ctx context.Context, userID int64) (*model.Voucher, error) {
	resp, err := r.cln.GetByUserID(ctx, &userv1.GetByUserIDRequest{UserId: userID})
	if err != nil {
		return nil, fmt.Errorf("failed to get user's voucher: %w", err)
	}

	return FromVoucher(resp.Voucher), nil
}
