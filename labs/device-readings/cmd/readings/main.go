package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"device-readings/internal/api"
	"device-readings/internal/db"
	"device-readings/internal/env"

	// "device-readings/internal/db"
	"device-readings/internal/queue"
	"device-readings/internal/readings/models"
	"device-readings/internal/readings/repo"

	"github.com/samber/lo"
)

func main() {
	start()
}

func start() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	config := lo.Must(env.New[queue.KafkaConfig]())

	writer := queue.NewKafkaWriter[[]models.BatchReading](config.Broker, config.Topic)
	defer writer.Close()

	db := lo.Must(db.NewFromEnv())
	repo := lo.Must(repo.NewGorm(db))
	server := lo.Must(api.New(writer, repo))

	srvErr := make(chan error, 1)
	go func() {
		srvErr <- server.Start()
	}()

	select {
	case err := <-srvErr:
		quit(ctx, err)
	case <-ctx.Done():
		slog.Info("shutting down")
		stop()
	}
}

func quit(ctx context.Context, err error) {
	slog.ErrorContext(ctx, err.Error())
	os.Exit(1)
}
