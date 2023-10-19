package app

import (
	"context"
	"github.com/caarlos0/env/v9"
	handler "github.com/malkev1ch/observability/userservice/internal/handler/grpc"
	"github.com/malkev1ch/observability/userservice/internal/repository"
	"github.com/malkev1ch/observability/userservice/internal/service"
	"log/slog"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	userv1 "github.com/malkev1ch/observability/userservice/gen/user/v1"
)

type App struct {
	cfg        *config
	grpcServer *grpc.Server
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
	return a.runGRPCServer()
}

func (a *App) Stop() error {
	return a.Stop()
}

func (a *App) init(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initSlog,
		a.initConfig,
		a.initGRPCServer,
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

func (a *App) initGRPCServer(_ context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	reflection.Register(a.grpcServer)

	userRepository := repository.NewUser()
	userService := service.NewUser(userRepository)
	userHandler := handler.NewUser(userService)

	userv1.RegisterUserServiceServer(a.grpcServer, userHandler)

	return nil
}

func (a *App) runGRPCServer() error {
	slog.Info(
		"GRPC server is started",
		slog.String("addr", a.cfg.Address),
	)

	list, err := net.Listen("tcp", a.cfg.Address)
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}
