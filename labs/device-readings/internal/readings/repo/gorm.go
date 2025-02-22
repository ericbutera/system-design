package repo

import (
	"context"
	"errors"

	"device-readings/internal/readings/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Gorm struct {
	db *gorm.DB
}

func NewGorm(db *gorm.DB) (*Gorm, error) {
	return &Gorm{
		db: db,
	}, nil
}

var ErrBatchReadingSaveFailure = errors.New("batch reading save failure")

func (g *Gorm) StoreReadings(readings []models.BatchReading) (StoreReadingsResult, error) {
	if len(readings) == 0 {
		return StoreReadingsResult{}, nil
	}

	var storedReadings []models.Reading
	for _, reading := range readings {
		storedReadings = append(storedReadings, models.Reading{
			DeviceID:    reading.DeviceID,
			ReadingType: reading.ReadingType,
			Timestamp:   reading.Timestamp,
			Value:       float64(reading.Value),
		})
	}

	result := StoreReadingsResult{}
	err := g.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&storedReadings)
		if res.Error != nil {
			return res.Error
		}

		result.Succeed = int(res.RowsAffected)
		return nil
	})

	return result, err
}

func (g *Gorm) GetReadings(ctx context.Context, filters Filters) ([]models.Reading, error) {
	var readings []models.Reading
	res := g.db.Model(&models.Reading{}).Find(&readings)
	if res.Error != nil {
		return nil, res.Error
	}
	return readings, nil
}
