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
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var config api.Config
	lo.Must0(env.Parse(&config))

	repo := lo.Must(db.New(config.DSN))

	svc := api.New(config, repo)
	server := svc.Server()

	srvErr := make(chan error, 1)
	go func() { srvErr <- server.Run() }()

	select {
	case err := <-srvErr:
		slog.Error("server error", "error", err)
		os.Exit(1)
	case <-ctx.Done():
		slog.Info("shutting down")
	}
}
