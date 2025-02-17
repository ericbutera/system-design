package queue

import (
	"context"
	"device-readings/internal/readings/models"
	"log/slog"
)

type KafkaProducer struct {
}

func NewKafkaProducer() *KafkaProducer {
	return &KafkaProducer{}
}

/*
	// Producer:
	broker := os.Getenv("KAFKA_BROKER")
	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker},
		Topic:   "readings",
	})
	defer producer.Close()

	ev := Readings{[]Reading{DeviceID:"device-123", Value:100}}
	data, _ := json.Marshal(ev)
	err = producer.WriteMessages(context.Background(), kafka.Message{ Value: data })
	if err != nil {
		log.Printf("Error producing message: %v", err)
	}
	fmt.Println("enqueued readings:", string(data))
*/

func (p *KafkaProducer) Write(ctx context.Context, readings []models.BatchReading) (WriteResult, error) {
	slog.Info("writing readings", "readings", readings)
	return WriteResult{ID: "todo-kafka-message-id"}, nil
}

type KafkaConsumer struct {
}

func NewKafkaConsumer() *KafkaConsumer {
	return &KafkaConsumer{}
}

/*
	// Consumer:
	broker := os.Getenv("KAFKA_BROKER")
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    "readings",
		GroupID:  "group-id",
		MinBytes: 10e3, // Minimum 10KB to fetch
		MaxBytes: 10e6, // Maximum 10MB to fetch
	})
	defer consumer.Close()

	for {
		msg, err := consumer.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}
		var data Readings
		err = json.Unmarshal(msg.Value, &data)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}
		fmt.Printf("Received readings: %+v\n", data)
	}
}
*/

func (c *KafkaConsumer) Read(ctx context.Context, handler func(ctx context.Context, readings []models.BatchReading) error) error {
	slog.Info("reading messages")
	return nil
}
