package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v11"
	"github.com/ericbutera/system-design/labs/url-shortener/internal/api"
	"github.com/ericbutera/system-design/labs/url-shortener/internal/db"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	eg, ctx := errgroup.WithContext(ctx)

	var config api.Config
	lo.Must0(env.Parse(&config))

	repo := lo.Must(db.New(config.DSN))

	svc := api.New(config, repo)
	server := svc.Server()

	eg.Go(func() error {
		return server.Run()
	})

	eg.Go(func() error {
		<-ctx.Done()
		slog.Info("shutting down")
		return ctx.Err()
	})

	if err := eg.Wait(); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
