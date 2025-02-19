package models

// TODO: use separate models for transport, domain logic and persistence layers
// TODO: reading data could be represented in a timeseries

import (
	"time"
)

// https://docs.timescale.com/quick-start/latest/golang/
type Device struct {
	ID       int64  `gorm:"primaryKey"                    json:"id"`
	DeviceID string `gorm:"uniqueIndex;type:varchar(255)" json:"device_id"`
	Type     string `gorm:"type:varchar(50)"              json:"type"`
	Location string `gorm:"type:varchar(50)"              json:"location"`
}

type Reading struct {
	Timestamp   time.Time `gorm:"not null;index:idx_unique_reading,unique"                   json:"timestamp"`
	DeviceID    string    `gorm:"type:varchar(255);not null;index:idx_unique_reading,unique" json:"device_id"`
	ReadingType string    `gorm:"type:varchar(50);not null;index:idx_unique_reading,unique"  json:"reading_type"`
	Value       float64   `gorm:"not null"                                                   json:"value"`
}

type BatchReading struct {
	DeviceID    string    `json:"device_id" binding:"required" description:"Device ID"                            example:"device-1"                  json:"device_id"`
	ReadingType string    `json:"type"      binding:"required" json:"reading_type"`
	Timestamp   time.Time `json:"timestamp" binding:"required" description:"Device reported timestamp of reading" example:"2021-01-01T00:00:00-05:00" json:"timestamp"`
	Value       float32   `json:"value"     binding:"required" description:"Reading data"                         example:"17"                        json:"count"`
}
