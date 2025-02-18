package main

import (
	"context"
	"device-readings/internal/env"
	"device-readings/internal/queue"
	"device-readings/internal/readings/models"
	"device-readings/internal/readings/processor"
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

	reader := queue.NewKafkaReader[[]models.BatchReading](config.Broker, config.Topic, config.Group)
	defer reader.Close()

	slog.Info("starting batch readings worker", "broker", config.Broker, "topic", config.Topic, "group", config.Group)

	processor := processor.NewProcessor(reader)
	err := processor.Run(ctx)

	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
