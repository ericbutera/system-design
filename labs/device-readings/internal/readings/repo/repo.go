package repo

import (
	"context"
	"time"

	"device-readings/internal/readings/models"
)

type RelativeTimes string

const (
	RelativeTimeQuarterHour RelativeTimes = "15m"
	RelativeTimeHalfHour    RelativeTimes = "30m"
	RelativeTimeOneHour     RelativeTimes = "1h"
	RelativeTimeOneDay      RelativeTimes = "1d"
)

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
}
