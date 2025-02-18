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
	slog.Info("running processor")
	return p.reader.Read(ctx, p.Handler)
}

func (p *Processor) Handler(ctx context.Context, batch []models.BatchReading) error {
	slog.Info("received readings", "readings", batch)
	res, err := p.repo.StoreReadings(batch)
	if err != nil {
		slog.Error("processor error storing readings", "error", err)
		return err
	}
	slog.Info("processor stored readings", "result", res)
	return nil
}
