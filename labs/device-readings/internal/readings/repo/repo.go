package repo

import (
	"context"
	"device-readings/internal/readings/models"
	"time"
)

type RelativeTimes string

const RelativeTimeQuarterHour RelativeTimes = "15m"
const RelativeTimeHalfHour RelativeTimes = "30m"
const RelativeTimeOneHour RelativeTimes = "1h"
const RelativeTimeOneDay RelativeTimes = "1d"

type Filters struct {
	DeviceID string
	StartAt  time.Time
	EndAt    time.Time
	Relative RelativeTimes
}

type StoreReadingsResult struct {
	Failures int
	Succeed  int
	ResultID string // queue message id
}

type Repo interface {
	StoreReadings(readings []models.BatchReading) (StoreReadingsResult, error)
	GetReadings(ctx context.Context, filters Filters) ([]models.Reading, error)
	GetReadingsByDevice(ctx context.Context, deviceID string) ([]models.Reading, error)
}
