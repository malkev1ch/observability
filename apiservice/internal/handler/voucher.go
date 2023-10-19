package handler

import (
	"context"
	"github.com/labstack/echo/v4"
	genv1 "github.com/malkev1ch/observability/apiservice/gen/v1"
	"github.com/malkev1ch/observability/apiservice/internal/model"
	"net/http"
)

type VoucherService interface {
	Search(ctx context.Context, userId int64) (*model.Voucher, error)
}

type Voucher struct {
	svc VoucherService
}

func NewVoucher(svc VoucherService) *Voucher {
	return &Voucher{svc: svc}
}

func (h *Voucher) SearchVoucher(ctx echo.Context, params genv1.SearchVoucherParams) error {
	voucher, err := h.svc.Search(ctx.Request().Context(), params.UserId)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return ctx.JSON(http.StatusOK, genv1.Voucher{
		CreatedAt: voucher.CreatedAt,
		Id:        voucher.ID,
		UserId:    voucher.UserId,
		UserName:  voucher.UserName,
		Value:     voucher.Value,
	})
}
