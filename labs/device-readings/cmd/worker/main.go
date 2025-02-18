package main

import (
	"context"
	"device-readings/internal/env"
	"device-readings/internal/readings/models"
	"device-readings/internal/readings/queue"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/samber/lo"
)

func main() {
	start()
}

func start() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	config := lo.Must(env.New[queue.KafkaConfig]())

	consumer := queue.NewKafkaConsumer(config.Broker, config.Topic, config.Group)
	defer consumer.Close()

	slog.Info("starting consumer", "broker", config.Broker, "topic", config.Topic, "group", config.Group)
	err := consumer.Read(ctx, func(ctx context.Context, batch []models.BatchReading) error {
		slog.Info("received readings", "readings", batch)
		return nil
	})
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
