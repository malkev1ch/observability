package client

import (
	"context"
	"fmt"
	"github.com/malkev1ch/observability/apiservice/internal/model"
	v1 "github.com/malkev1ch/observability/voucherservice/gen/voucher/v1"
)

type Voucher struct {
	cln v1.VoucherServiceClient
}

func NewVoucher(cln v1.VoucherServiceClient) *Voucher {
	return &Voucher{cln: cln}
}

func FromVoucher(in *v1.Voucher) *model.Voucher {
	return &model.Voucher{
		ID:        in.Id,
		Value:     in.Value,
		UserId:    in.UserId,
		UserName:  "",
		CreatedAt: in.CreatedAt.AsTime(),
	}
}

func (r *Voucher) GetByUserID(ctx context.Context, userID int64) (*model.Voucher, error) {
	resp, err := r.cln.GetByUserID(ctx, &v1.GetByUserIDRequest{UserId: userID})
	if err != nil {
		return nil, fmt.Errorf("failed to get user's voucher: %w", err)
	}

	return FromVoucher(resp.Voucher), nil
}
