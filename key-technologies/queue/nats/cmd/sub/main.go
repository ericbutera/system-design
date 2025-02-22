package main

import (
	"log"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
)

func main() {
	url := "nats://nats:4222"
	nc := lo.Must(nats.Connect(url))
	js := lo.Must(nc.JetStream())

	// Create a durable consumer (subscription)
	_, err := js.Subscribe("events.*", func(msg *nats.Msg) {
		defer msg.Nak()
		slog.Info("received message", "msg.Header", msg.Header, "data", string(msg.Data))

		if err := msg.Ack(); err != nil {
			log.Fatalf("Error acknowledging message: %v", err)
		}
	}, nats.Durable("EVENTS-DURABLE"), nats.ManualAck())
	if err != nil {
		log.Fatalf("Error subscribing to messages: %v", err)
	}

	// Keep the connection alive
	select {}
}
