package handler

import (
	"context"
	"net/http"

	genv1 "github.com/malkev1ch/observability/services/apiservice/gen/v1"
	"github.com/malkev1ch/observability/services/apiservice/internal/model"

	"github.com/labstack/echo/v4"
)

type VoucherService interface {
	Search(ctx context.Context, userID int64) (*model.Voucher, error)
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
		UserId:    voucher.UserID,
		UserName:  voucher.UserName,
		Value:     voucher.Value,
	})
}
