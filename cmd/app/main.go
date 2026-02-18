package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"

	"go-users/internal/app"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	app, err := app.NewApp(ctx)

	if err != nil {
		panic(err)
	}

	go func() {
		if err := app.Start(ctx); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-ctx.Done()

	errors := app.Shutdown()

	for _, err := range errors {
		if err != context.Canceled {
		}
	}
}
