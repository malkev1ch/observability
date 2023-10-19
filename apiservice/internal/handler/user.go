package handler

import (
	"context"
	"github.com/labstack/echo/v4"
	genv1 "github.com/malkev1ch/observability/apiservice/gen/v1"
	"github.com/malkev1ch/observability/apiservice/internal/model"
	"net/http"
	"time"
)

type UserService interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
	Create(ctx context.Context, user *model.User) (*model.User, error)
}

type User struct {
	svc UserService
}

func NewUser(svc UserService) *User {
	return &User{svc: svc}
}

func (h *User) GetUserById(ctx echo.Context, id int64) error {
	user, err := h.svc.GetByID(ctx.Request().Context(), id)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return ctx.JSON(http.StatusOK, genv1.User{
		CreatedAt: user.CreatedAt,
		Id:        user.ID,
		Name:      user.Name,
	})
}

func (h *User) CreateUser(ctx echo.Context) error {
	var payload genv1.CreateUser
	err := ctx.Bind(&payload)
	if err != nil {
		return echo.ErrBadRequest
	}

	user, err := h.svc.Create(
		ctx.Request().Context(),
		&model.User{
			ID:        0,
			Name:      payload.Name,
			CreatedAt: time.Time{},
		},
	)

	return ctx.JSON(http.StatusOK, genv1.User{
		CreatedAt: user.CreatedAt,
		Id:        user.ID,
		Name:      user.Name,
	})
}
