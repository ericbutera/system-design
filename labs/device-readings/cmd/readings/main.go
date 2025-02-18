package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"device-readings/internal/api"
	"device-readings/internal/env"

	// "device-readings/internal/db"
	"device-readings/internal/readings/queue"
	"device-readings/internal/readings/repo"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

func main() {
	start()
}

func start() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	config := lo.Must(env.New[queue.KafkaConfig]())

	producer := queue.NewKafkaProducer(config.Broker, config.Topic)
	defer producer.Close()

	db := &gorm.DB{} // db := lo.Must(db.New(db.NewDefaultConfig()))
	repo := repo.NewGorm(db)
	server := lo.Must(api.New(producer, repo))

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
