package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"device-readings/internal/db"
	"device-readings/internal/env"
	"device-readings/internal/queue"
	"device-readings/internal/readings/models"
	"device-readings/internal/readings/processor"
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

	reader := queue.NewKafkaReader[[]models.BatchReading](config.Broker, config.Topic, config.Group)
	defer reader.Close()

	slog.Info("starting batch readings worker", "broker", config.Broker, "topic", config.Topic, "group", config.Group)

	db := lo.Must(db.NewFromEnv())
	repo := lo.Must(repo.NewGorm(db))

	processor := processor.NewProcessor(reader, repo)
	err := processor.Run(ctx)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
