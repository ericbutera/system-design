package queue

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

const (
	MinBytes = 100
	MaxBytes = 10e6 // Maximum 10MB to fetch
)

type KafkaConfig struct {
	Broker string `env:"KAFKA_BROKER" envDefault:"redpanda:9092"`
	Topic  string `env:"KAFKA_TOPIC"  envDefault:"readings-topic"`
	Group  string `env:"KAFKA_GROUP"  envDefault:"readings-group"`
}

type KafkaWriter[T any] struct {
	writer *kafka.Writer
}

func NewKafkaWriter[T any](broker string, topic string) *KafkaWriter[T] {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(broker),
		Topic:        topic,
		RequiredAcks: kafka.RequireAll,
	}

	return &KafkaWriter[T]{
		writer: writer,
	}
}

func (p *KafkaWriter[T]) Close() {
	p.writer.Close()
}

func (p *KafkaWriter[T]) Write(ctx context.Context, data T) error {
	encoded, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		return err
	}

	err = p.writer.WriteMessages(context.Background(), kafka.Message{Value: encoded})
	if err != nil {
		log.Printf("Error producing message: %v", err)
		return err
	}

	return nil
}

type KafkaReader[T any] struct {
	reader *kafka.Reader
}

func NewKafkaReader[T any](broker, topic, group string) *KafkaReader[T] {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  group, // group stores the offset; group is specific to a worker/subscriber/consumer
		MinBytes: MinBytes,
		MaxBytes: MaxBytes,
		// StartOffset: kafka.FirstOffset, // this can be used to reprocess messages
		// CommitInterval: 0,
		// TODO no auto-commit! only ack when side-effects are done
	})
	return &KafkaReader[T]{
		reader: reader,
	}
}

func (c *KafkaReader[T]) Close() {
	c.reader.Close()
}

// TODO: dead-letter queue
// Read reads messages from Kafka and calls the handler function for each message.
// Return false in the handler to retry the message.
func (c *KafkaReader[T]) Read(ctx context.Context, handler func(ctx context.Context, data T) error) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		slog.Info("reading message",
			"headers", msg.Headers,
			"high watermark", msg.HighWaterMark,
			"offset", msg.Offset,
			"key", msg.Key,
			"value", string(msg.Value),
		)
		if err != nil {
			slog.Error("error reading message", "error", err)
			return err
		}

		var data T
		err = json.Unmarshal(msg.Value, &data) // TODO: commit message if unmarshal fails to not reprocess the same message
		if err != nil {
			return err
		}

		err = handler(ctx, data)
		if err != nil {
			return err
		}

		err = c.reader.CommitMessages(ctx, msg)
		if err != nil {
			return err
		}
		slog.Info("message committed successfully")
	}
}
