package client

import (
	"context"
	"fmt"

	v1 "github.com/malkev1ch/observability/voucherservice/gen/voucher/v1"
)

type Voucher struct {
	cln v1.VoucherServiceClient
}

func NewVoucher(cln v1.VoucherServiceClient) *Voucher {
	return &Voucher{cln: cln}
}

func (r *Voucher) Create(ctx context.Context, userID int64) error {
	_, err := r.cln.Create(ctx, &v1.CreateRequest{UserId: userID})
	if err != nil {
		return fmt.Errorf("failed to create voucher: %w", err)
	}

	return nil
}
