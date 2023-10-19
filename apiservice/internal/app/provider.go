package app

import (
	"github.com/malkev1ch/observability/apiservice/internal/handler"
	"github.com/malkev1ch/observability/apiservice/internal/repository/client"
	"github.com/malkev1ch/observability/apiservice/internal/service"
	userv1 "github.com/malkev1ch/observability/userservice/gen/user/v1"
	"google.golang.org/grpc"
)

type provider struct {
	Handler *handler.Handler

	voucherHandler *handler.Voucher

	userHandler    *handler.User
	userService    *service.User
	userRepository *client.User
	userConn       *grpc.ClientConn
	userClient     userv1.UserServiceClient
}
