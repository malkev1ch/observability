package app

import (
	"context"
	"github.com/caarlos0/env/v9"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	genv1 "github.com/malkev1ch/observability/apiservice/gen/v1"
	"github.com/malkev1ch/observability/apiservice/internal/handler"
	"github.com/malkev1ch/observability/apiservice/internal/repository/client"
	"github.com/malkev1ch/observability/apiservice/internal/service"
	userv1 "github.com/malkev1ch/observability/userservice/gen/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"os"
	"time"
)

type App struct {
	cfg        *config
	provider   *provider
	echoServer *echo.Echo
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.init(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	return a.runEchoServer()
}

func (a *App) Stop() error {
	err := a.Stop()
	if err != nil {
		return err
	}

	err = a.provider.userConn.Close()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) init(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initSlog,
		a.initConfig,
		a.initProvider,
		a.initEchoServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		return err
	}

	a.cfg = &cfg
	return nil
}

func (a *App) initSlog(_ context.Context) error {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Value.Kind() == slog.KindTime {
				a.Value = slog.StringValue(a.Value.Time().In(time.UTC).Format(time.RFC3339))
			}

			return a
		},
	})))

	return nil
}

func (a *App) initProvider(_ context.Context) error {
	grpcDialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	userConn, err := grpc.Dial(a.cfg.UserServiceAddress, grpcDialOpts...)
	if err != nil {
		return nil
	}

	userClient := userv1.NewUserServiceClient(userConn)
	userRepository := client.NewUser(userClient)
	userService := service.NewUser(userRepository)
	userHandler := handler.NewUser(userService)

	voucherHandler := handler.NewVoucher(nil)

	a.provider = &provider{
		Handler:        handler.New(userHandler, voucherHandler),
		voucherHandler: voucherHandler,
		userHandler:    userHandler,
		userService:    userService,
		userRepository: userRepository,
		userConn:       userConn,
		userClient:     userClient,
	}

	return nil
}

func (a *App) initEchoServer(_ context.Context) error {
	a.echoServer = echo.New()

	a.echoServer.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	a.echoServer.HideBanner = true
	a.echoServer.HidePort = true

	genv1.RegisterHandlers(a.echoServer, a.provider.Handler)

	routes := a.echoServer.Routes()

	for _, route := range routes {
		slog.Info(
			"Registered route",
			slog.String("path", route.Path),
			slog.String("method", route.Method),
		)
	}

	return nil
}

func (a *App) runEchoServer() error {
	slog.Info(
		"Echo server is started",
		slog.String("addr", a.cfg.Address),
	)

	err := a.echoServer.Start(a.cfg.Address)
	if err != nil {
		return err
	}

	return nil
}
