package processor

import (
	"context"
	"device-readings/internal/queue"
	"device-readings/internal/readings/models"
	"device-readings/internal/readings/repo"
	"log/slog"
)

type ReaderType = queue.KafkaReader[[]models.BatchReading]

// TODO: this should be called "BatchReadingProcessor" - make a generic message handler
type Processor struct {
	reader *ReaderType
	repo   repo.Repo
}

func NewProcessor(reader *ReaderType, repo repo.Repo) *Processor {
	return &Processor{
		reader: reader,
		repo:   repo,
	}
}

func (p *Processor) Run(ctx context.Context) error {
	return p.reader.Read(ctx, p.Handler)
}

func (p *Processor) Handler(ctx context.Context, batch []models.BatchReading) error {
	slog.Debug("received readings", "count", len(batch))
	res, err := p.repo.StoreReadings(batch)
	if err != nil {
		return err
	}
	slog.Debug("processor stored readings", "result", res)
	return nil
}
