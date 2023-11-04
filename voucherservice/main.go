package main

import (
	"context"
	"github.com/caarlos0/env/v9"
	voucherv1 "github.com/malkev1ch/observability/voucherservice/gen/voucher/v1"
	handler "github.com/malkev1ch/observability/voucherservice/internal/handler/grpc"
	"github.com/malkev1ch/observability/voucherservice/internal/model"
	"github.com/malkev1ch/observability/voucherservice/internal/repository"
	"github.com/malkev1ch/observability/voucherservice/internal/service"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const appName = "voucher-service"

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

	// Establish connection to opentelemetry agent
	conn, err := grpc.DialContext(ctx, cfg.OtelAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Failed to init opentelemetry conn", slog.String("error", err.Error()))
		return
	}

	//Initialize an exporter
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

	server := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)
	reflection.Register(server)

	voucherRepository := repository.NewVoucher()
	voucherService := service.NewVoucher(voucherRepository)
	voucherHandler := handler.NewVoucher(voucherService)

	voucherv1.RegisterVoucherServiceServer(server, voucherHandler)

	list, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		slog.Error("Failed to init net.Listen", slog.String("error", err.Error()))
	}

	go func() {
		slog.Info("GRPC server is started", slog.String("addr", cfg.Address))
		err = server.Serve(list)
		if err != nil {
			slog.Error("Failed to start GRPC server", slog.String("error", err.Error()))
		}
		stop()
	}()

	<-ctx.Done()

	server.Stop()
	slog.Info("Successfully stopped app")

}
