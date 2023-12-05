package client

import (
	"context"
	"fmt"

	userv1 "github.com/malkev1ch/observability/services/voucherservice/gen/voucher/v1"
)

type Voucher struct {
	cln userv1.VoucherServiceClient
}

func NewVoucher(cln userv1.VoucherServiceClient) *Voucher {
	return &Voucher{cln: cln}
}

func (r *Voucher) Create(ctx context.Context, userID int64) error {
	_, err := r.cln.Create(ctx, &userv1.CreateRequest{UserId: userID})
	if err != nil {
		return fmt.Errorf("failed to create voucher: %w", err)
	}

	return nil
}
