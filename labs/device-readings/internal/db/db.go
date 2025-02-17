package db

import (
	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config struct {
	LogQueries bool
	AutoModels []any // Intended to be gorm.Model{}
}

func NewDefaultConfig() *Config {
	return &Config{
		LogQueries: true,
	}
}

func New(config *Config) (*gorm.DB, error) {
	opts := &gorm.Config{}
	if config.LogQueries {
		opts.Logger = slogGorm.New()
	}

	driver := sqlite.Open("file::memory:?cache=shared") // TODO: configure driver
	instance, err := gorm.Open(driver, opts)
	if err != nil {
		return nil, err
	}

	if config.AutoModels != nil {
		if err := instance.AutoMigrate(config.AutoModels...); err != nil {
			return nil, err
		}
	}

	return instance, nil
}
