package main

import (
	"context"
	"github.com/malkev1ch/observability/userservice/internal/app"
	"log/slog"
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

	err = a.Run()
	if err != nil {
		slog.Error(
			"Failed to run app",
			slog.String("service", appName),
			slog.String("error", err.Error()),
		)
		return
	}
}
