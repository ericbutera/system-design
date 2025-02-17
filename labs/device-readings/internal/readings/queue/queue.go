package queue

import (
	"context"
	"device-readings/internal/readings/models"
)

type WriteResult struct {
	ID string
}

type Producer interface {
	Write(ctx context.Context, batch []models.BatchReading) (WriteResult, error) // TODO: generic type
}

type Consumer interface {
	Read(ctx context.Context, handler func(ctx context.Context, batch []models.BatchReading) error) error
}
