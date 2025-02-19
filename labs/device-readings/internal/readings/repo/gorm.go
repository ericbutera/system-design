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

func NewGorm(db *gorm.DB) (*Gorm, error) {
	if err := db.AutoMigrate(&models.Reading{}); err != nil {
		return nil, err
	}
	return &Gorm{
		db: db,
	}, nil
}

var (
	ErrBatchReadingSaveFailure = errors.New("batch reading save failure")
)

func (g *Gorm) StoreReadings(readings []models.BatchReading) (StoreReadingsResult, error) {
	// TODO: save failures for inspection
	result := StoreReadingsResult{}
	err := g.db.Transaction(func(tx *gorm.DB) error {
		for _, reading := range readings {
			res := tx.Create(models.Reading{
				DeviceID:    reading.DeviceID,
				ReadingType: reading.ReadingType,
				Timestamp:   reading.Timestamp,
				Value:       reading.Value,
			})
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
	var readings []models.Reading
	res := g.db.Model(&models.Reading{}).Find(&readings)
	if res.Error != nil {
		return nil, res.Error
	}
	return readings, nil
}

func (g *Gorm) GetReadingsByDevice(ctx context.Context, deviceID string) ([]models.Reading, error) {
	slog.Info("getting readings by device", "deviceID", deviceID)
	return []models.Reading{}, nil
}
