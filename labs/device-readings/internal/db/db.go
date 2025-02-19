package db

import (
	"log/slog"

	"device-readings/internal/env"
	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	LogQueries bool   `env:"DB_LOG_QUERIES" envDefault:"true"`
	Dsn        string `env:"DB_DSN"         envDefault:"host=timescaledb user=postgres password=password dbname=postgres port=5432 sslmode=disable"`
}

func NewFromEnv() (*gorm.DB, error) {
	config, err := env.New[Config]()
	if err != nil {
		return nil, err
	}
	return New(config)
}

func New(config *Config) (*gorm.DB, error) {
	opts := &gorm.Config{}
	if config.LogQueries {
		l := slogGorm.New()
		slogGorm.SetLogLevel(slogGorm.DefaultLogType, slog.LevelDebug)
		opts.Logger = l
	}

	instance, err := gorm.Open(postgres.Open(config.Dsn), opts)
	if err != nil {
		return nil, err
	}

	return instance, nil
}
