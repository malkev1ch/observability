package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	genv1 "github.com/malkev1ch/observability/services/apiservice/gen/v1"
	"github.com/malkev1ch/observability/services/apiservice/internal/handler"
	"github.com/malkev1ch/observability/services/apiservice/internal/model"
	"github.com/malkev1ch/observability/services/apiservice/internal/repository/client"
	"github.com/malkev1ch/observability/services/apiservice/internal/service"
	userv1 "github.com/malkev1ch/observability/services/userservice/gen/user/v1"
	voucherv1 "github.com/malkev1ch/observability/services/voucherservice/gen/voucher/v1"

	"github.com/caarlos0/env/v9"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	swaggergui "github.com/malkev1ch/observability/pkg/swaggergui"
)

const appName = "api-service"

//nolint:funlen
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Value.Kind() == slog.KindTime {
				a.Value = slog.StringValue(a.Value.Time().In(time.UTC).Format(time.RFC3339))
			}
			return a
		},
	}).WithAttrs([]slog.Attr{{Key: "service", Value: slog.StringValue(appName)}})))

	var cfg model.Config
	err := env.Parse(&cfg)
	if err != nil {
		slog.Error("Failed to init config", slog.String("error", err.Error()))
		return
	}

	// ------------------ OPENTELEMETRY CONFIGURATION ------------------ //

	// Identify application as resource
	r, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceVersion("v0.1"),
			semconv.ServiceName(appName),
		),
	)
	if err != nil {
		slog.Error("Failed to create opentelemetry resource", slog.String("error", err.Error()))
		return
	}

	// Establish connection to opentelemetry agent
	conn, err := grpc.DialContext(ctx, cfg.OtelAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Failed to init opentelemetry conn", slog.String("error", err.Error()))
		return
	}

	// Initialize an exporter
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		slog.Error("Failed to init opentelemetry exporter", slog.String("error", err.Error()))
		return
	}

	// Initialize a new tracer provider with a batch span processor and the given exporter.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(r),
		// Samples every trace
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		// Can be configured
		sdktrace.WithBatcher(exporter),
	)

	// Set as a global provider
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// -------------- END OPENTELEMETRY CONFIGURATION ------------------ //

	grpcDialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}

	userConn, err := grpc.Dial(cfg.UserServiceAddress, grpcDialOpts...)
	if err != nil {
		slog.Error("Failed to init conn to user service", slog.String("error", err.Error()))
		return
	}

	voucherConn, err := grpc.Dial(cfg.VoucherServiceAddress, grpcDialOpts...)
	if err != nil {
		slog.Error("Failed to init conn to voucher service", slog.String("error", err.Error()))
		return
	}

	voucherClient := voucherv1.NewVoucherServiceClient(voucherConn)
	userClient := userv1.NewUserServiceClient(userConn)

	voucherRepository := client.NewVoucher(voucherClient)
	userRepository := client.NewUser(userClient)

	userService := service.NewUser(userRepository)
	voucherService := service.NewVoucher(voucherRepository, userRepository)

	userHandler := handler.NewUser(userService)
	voucherHandler := handler.NewVoucher(voucherService)

	e := echo.New()
	e.Use(slogecho.NewWithFilters(slog.Default(), slogecho.IgnorePathContains("healthz")))
	e.Use(middleware.Recover())
	e.Use(otelecho.Middleware("apiservice", otelecho.WithSkipper(func(c echo.Context) bool {
		return strings.Contains(c.Request().URL.Path, "healthz")
	}),
	))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	e.HideBanner = true
	e.HidePort = true

	swaggerGuiHandler := swaggergui.NewHandlerWithConfig(swaggergui.Config{
		Title:              "Observability",
		SwaggerJSON:        "/api/v1/openapi.json",
		BasePath:           "/v1/docs/",
		ShowTopBar:         false,
		HideCurl:           false,
		JSONEditor:         false,
		PreAuthorizeAPIKey: nil,
		SettingsUI:         nil,
	})

	e.GET("/v1/docs", func(c echo.Context) error {
		swaggerGuiHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})
	e.GET("/healthz", func(c echo.Context) error { return c.NoContent(http.StatusOK) })

	genv1.RegisterHandlers(e, handler.New(userHandler, voucherHandler))

	for _, route := range e.Routes() {
		slog.Info(
			"Registered route",
			slog.String("path", route.Path),
			slog.String("method", route.Method),
		)
	}

	go func() {
		slog.Info("Echo server is started", slog.String("addr", cfg.Address))
		err = e.Start(cfg.Address)
		if err != nil {
			slog.Error("Failed to start echo server", slog.String("error", err.Error()))
		}
		stop()
	}()

	<-ctx.Done()
	err = e.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to stop app", slog.String("error", err.Error()))
		return
	}

	slog.Info("Successfully stopped app")
}
