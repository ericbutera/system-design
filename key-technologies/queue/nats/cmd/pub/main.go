package main

import (
	"log"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
)

func main() {
	url := "nats://nats:4222"
	nc := lo.Must(nats.Connect(url))
	js := lo.Must(nc.JetStream())

	// Create a stream (equivalent to a topic in GCP Pub/Sub)
	_, err := js.AddStream(&nats.StreamConfig{
		Name:     "EVENTS",
		Subjects: []string{"events.*"},
	})
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	for {
		result, err := js.Publish("events.created", []byte("Hello, JetStream!!"))
		if err != nil {
			log.Fatalf("Error publishing message: %v", err)
		}
		slog.Info("Published message", "result", result)
		time.Sleep(1 * time.Minute)
	}
}
