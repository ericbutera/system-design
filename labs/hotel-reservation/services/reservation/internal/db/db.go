package db

import (
	"time"

	gorm_logrus "github.com/onrik/gorm-logrus"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	LogQueries  bool
	AutoMigrate bool
}

func NewDefaultConfig() *Config {
	return &Config{
		LogQueries:  true,
		AutoMigrate: true,
	}
}

func New(config *Config) (*gorm.DB, error) {
	opts := &gorm.Config{}
	if config.LogQueries { // TODO: convert to slog
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetLevel(logrus.DebugLevel)
		opts.Logger = gorm_logrus.New()
	}

	dsn := "host=pg-postgresql user=postgres password=password dbname=postgres port=5432 sslmode=disable" // TODO: env
	d, err := gorm.Open(postgres.Open(dsn), opts)
	if err != nil {
		return nil, err
	}

	d.NowFunc = func() time.Time {
		return time.Now().UTC()
	}

	return d, nil
}
