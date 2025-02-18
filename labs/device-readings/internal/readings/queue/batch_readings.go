package queue

import (
	"context"
	"device-readings/internal/readings/models"
)

type BatchReadingWriter interface {
	Write(ctx context.Context, batch []models.BatchReading) error
}

type BatchReadingReader interface {
	Read(ctx context.Context, handler func(ctx context.Context, batch []models.BatchReading) error) error
}
