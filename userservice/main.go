package main

import (
	"context"
	"github.com/malkev1ch/observability/userservice/internal/app"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const appName = "user-service"

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		slog.Error(
			"Failed to init app",
			slog.String("service", appName),
			slog.String("error", err.Error()),
		)
		return
	}

	go func() {
		err = a.Run()
		if err != nil {
			slog.Error(
				"Failed to run app",
				slog.String("service", appName),
				slog.String("error", err.Error()),
			)
			return
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	<-done

	err = a.Stop()
	if err != nil {
		slog.Error(
			"Failed to stop app",
			slog.String("service", appName),
			slog.String("error", err.Error()),
		)
		return
	}

	slog.Info(
		"Successfully stopped app",
		slog.String("service", appName),
	)
}
