package handler

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "github.com/malkev1ch/observability/voucherservice/gen/voucher/v1"
	"github.com/malkev1ch/observability/voucherservice/internal/model"
)

type VoucherService interface {
	GetByUserID(ctx context.Context, id int64) (*model.Voucher, error)
	Create(ctx context.Context, voucher *model.Voucher) (*model.Voucher, error)
}

type Voucher struct {
	svc VoucherService

	v1.UnimplementedVoucherServiceServer
}

func NewVoucher(svc VoucherService) *Voucher {
	return &Voucher{svc: svc}
}

func ToVoucher(voucher *model.Voucher) *v1.Voucher {
	return &v1.Voucher{
		Id:        voucher.ID,
		UserId:    voucher.UserID,
		Value:     voucher.Value,
		CreatedAt: timestamppb.New(voucher.CreatedAt),
	}
}

func (h *Voucher) Create(ctx context.Context, in *v1.CreateRequest) (*v1.CreateResponse, error) {
	voucher, err := h.svc.Create(ctx, &model.Voucher{
		ID:        0,
		UserID:    in.UserId,
		Value:     "",
		CreatedAt: time.Time{},
	})
	if err != nil {
		slog.Error(
			"Failed to create voucher",
			slog.Any("request", in),
		)
		return nil, err
	}

	return &v1.CreateResponse{Voucher: ToVoucher(voucher)}, nil
}

func (h *Voucher) GetByUserID(ctx context.Context, in *v1.GetByUserIDRequest) (*v1.GetByUserIDResponse, error) {
	voucher, err := h.svc.GetByUserID(ctx, in.UserId)
	if err != nil {
		slog.Error(
			"Failed to get voucher by user id",
			slog.Any("request", in),
		)
		return nil, err
	}

	return &v1.GetByUserIDResponse{Voucher: ToVoucher(voucher)}, nil
}
