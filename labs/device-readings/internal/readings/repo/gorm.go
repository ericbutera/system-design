package repo

import (
	"context"
	"device-readings/internal/readings/models"
	"errors"
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

var (
	ErrBatchReadingSaveFailure = errors.New("batch reading save failure")
)

func (g *Gorm) StoreReadings(readings []models.BatchReading) (StoreReadingsResult, error) {
	// TODO: save failures for inspection
	result := StoreReadingsResult{}
	err := g.db.Transaction(func(tx *gorm.DB) error {
		for _, reading := range readings {
			res := tx.Create(&reading)
			if res.Error != nil {
				result.Failures++
				continue
			}
			result.Succeed++
		}
		if result.Failures > 0 {
			tx.Rollback() // TODO: determine failure threshold before rollback
			return ErrBatchReadingSaveFailure
		}
		return nil
	})
	return result, err
}

func (g *Gorm) GetReadings(ctx context.Context, filters Filters) ([]models.Reading, error) {
	slog.Info("getting readings", "filters", filters)
	return []models.Reading{}, nil
}

func (g *Gorm) GetReadingsByDevice(ctx context.Context, deviceID string) ([]models.Reading, error) {
	slog.Info("getting readings by device", "deviceID", deviceID)
	return []models.Reading{}, nil
}
