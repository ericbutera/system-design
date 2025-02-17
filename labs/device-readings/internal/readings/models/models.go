package models

// TODO: use separate models for transport, domain logic and persistence layers
// TODO: reading data could be represented in a timeseries

import (
	"time"
)

type Reading struct {
	DeviceID    string    `binding:"required"       description:"Device ID"                            example:"device-1"                  json:"device_id"`
	ReadingType string    `binding:"required"       description:"Type of reading"                      example:"temperature"                json:"reading_type"`
	Timestamp   time.Time `binding:"required"       description:"Device reported timestamp of reading" example:"2021-01-01T00:00:00-05:00" json:"timestamp"`
	Value       float32   `binding:"required,min=0" description:"Reading data"                         example:"17"                        json:"count"`
	// TODO: compare device reported Timestamp to server CreatedAt to determine time integrity; DDIA: unreliable clocks
}

type BatchReading struct {
	DeviceID    string    `json:"device_id" binding:"required" description:"Device ID"                       example:"device-1"                  json:"device_id"`
	ReadingType string    `json:"type" binding:"required" description:"Type of reading"                      	example:"temperature"                json:"reading_type"`
	Timestamp   time.Time `json:"timestamp" binding:"required" description:"Device reported timestamp of reading" example:"2021-01-01T00:00:00-05:00" json:"timestamp"`
	Value       float32   `json:"value" binding:"required" description:"Reading data"                         example:"17"                        json:"count"`
}
