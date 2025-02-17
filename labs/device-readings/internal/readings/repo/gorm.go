package repo

import (
	"context"
	"device-readings/internal/readings/models"
	"log/slog"

	"gorm.io/gorm"
)

type Gorm struct {
	db *gorm.DB
}

func NewGorm(db *gorm.DB) *Gorm {
	return &Gorm{
		db: db,
	}
}

func (g *Gorm) StoreReadings(readings []models.Reading) (StoreReadingsResult, error) {
	slog.Debug("storing readings", "readings", readings)
	return StoreReadingsResult{ResultID: "result-id"}, nil
}

func (g *Gorm) GetReadings(ctx context.Context, filters Filters) ([]models.Reading, error) {
	slog.Info("getting readings", "filters", filters)
	return []models.Reading{}, nil
}

func (g *Gorm) GetReadingsByDevice(ctx context.Context, deviceID string) ([]models.Reading, error) {
	slog.Info("getting readings by device", "deviceID", deviceID)
	return []models.Reading{}, nil
}
