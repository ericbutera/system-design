package queue

import (
	"context"
	"device-readings/internal/readings/models"
	"encoding/json"
	"log"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

type KafkaConfig struct {
	Broker string `env:"KAFKA_BROKER" envDefault:"redpanda:9092"`
	Topic  string `env:"KAFKA_TOPIC" envDefault:"readings"`
	Group  string `env:"KAFKA_GROUP" envDefault:"readings-group"`
}

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(broker string, topic string) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(broker),
		Topic:        topic,
		RequiredAcks: kafka.RequireAll,
	}

	return &KafkaProducer{
		writer: writer,
	}
}

func (p *KafkaProducer) Close() {
	p.writer.Close()
}

func (p *KafkaProducer) Write(ctx context.Context, readings []models.BatchReading) error {
	slog.Info("writing readings", "readings", readings)
	data, err := json.Marshal(readings)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		return err
	}

	err = p.writer.WriteMessages(context.Background(), kafka.Message{Value: data})
	if err != nil {
		log.Printf("Error producing message: %v", err)
		return err
	}

	return nil
}

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer(broker, topic, group string) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  group,
		MinBytes: 1,
		MaxBytes: 10e6, // Maximum 10MB to fetch
		//StartOffset: kafka.FirstOffset, // TODO: remove
		//CommitInterval: 0,
	})
	return &KafkaConsumer{
		reader: reader,
	}
}

func (c *KafkaConsumer) Close() {
	c.reader.Close()
}

func (c *KafkaConsumer) Read(ctx context.Context, handler func(ctx context.Context, readings []models.BatchReading) error) error {
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

		var data []models.BatchReading
		err = json.Unmarshal(msg.Value, &data)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}

		err = handler(ctx, data)
		if err != nil {
			slog.Error("error handling message", "error", err)
			return err
		}

		// err = c.reader.CommitMessages(ctx, msg)
		// if err != nil {
		// 	slog.Error("error committing message", "error", err)
		// 	return err
		// }
		slog.Info("message committed successfully")
	}
}
