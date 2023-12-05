package handler

import (
	"context"
	"log/slog"
	"time"

	userv1 "github.com/malkev1ch/observability/services/voucherservice/gen/voucher/v1"
	"github.com/malkev1ch/observability/services/voucherservice/internal/model"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type VoucherService interface {
	GetByUserID(ctx context.Context, id int64) (*model.Voucher, error)
	Create(ctx context.Context, voucher *model.Voucher) (*model.Voucher, error)
}

type Voucher struct {
	svc VoucherService

	userv1.UnimplementedVoucherServiceServer
}

func NewVoucher(svc VoucherService) *Voucher {
	return &Voucher{svc: svc}
}

func ToVoucher(voucher *model.Voucher) *userv1.Voucher {
	return &userv1.Voucher{
		Id:        voucher.ID,
		UserId:    voucher.UserID,
		Value:     voucher.Value,
		CreatedAt: timestamppb.New(voucher.CreatedAt),
	}
}

func (h *Voucher) Create(ctx context.Context, in *userv1.CreateRequest) (*userv1.CreateResponse, error) {
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

	return &userv1.CreateResponse{Voucher: ToVoucher(voucher)}, nil
}

func (h *Voucher) GetByUserID(ctx context.Context, in *userv1.GetByUserIDRequest) (*userv1.GetByUserIDResponse, error) {
	voucher, err := h.svc.GetByUserID(ctx, in.UserId)
	if err != nil {
		slog.Error(
			"Failed to get voucher by user id",
			slog.Any("request", in),
		)
		return nil, err
	}

	return &userv1.GetByUserIDResponse{Voucher: ToVoucher(voucher)}, nil
}
