package reporting

import (
	"context"
	"device-readings/internal/models"
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

// Reporting & Analytics service
type Reporting interface {
	GetReadings(filters Filters) ([]*models.Reading, error)
	GetReadingsByDevice(ctx context.Context, deviceID string) ([]models.Reading, error)
}
