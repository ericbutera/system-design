package db

import (
	"errors"
	"strings"

	"github.com/ericbutera/system-design/labs/url-shortener/internal/base62"
	"github.com/ericbutera/system-design/labs/url-shortener/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrDuplicate = errors.New("duplicate slug")
	ErrNotFound  = errors.New("not found")
)

type DB struct {
	db *gorm.DB
}

func New(dsn string) (*DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return &DB{db: db}, err
}

type SlugResult struct {
	Slug    string
	Counter int64
}

func (d *DB) GenerateSlug() (*SlugResult, error) {
	var counter int64
	res := d.db.Raw("SELECT nextval('url_counter')").Scan(&counter)
	if res.Error != nil {
		return nil, res.Error
	}
	return &SlugResult{
		Slug:    base62.EncodeInt64(counter),
		Counter: counter,
	}, nil
}

func (d *DB) CreateURL(url *models.URL) error {
	res := d.db.Create(&url)
	if res.Error != nil {
		if strings.Contains(res.Error.Error(), "violates unique constraint") {
			return ErrDuplicate
		}
	}
	return nil
}

func (d *DB) GetURL(slug string) (*models.URL, error) {
	var url models.URL
	res := d.db.Where("slug = ?", slug).First(&url)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, res.Error
	}
	return &url, nil
}

func (d *DB) GetStats(slug string) (*models.URLStats, error) {
	panic("not implemented")
	// var stats models.URLStats
	// res := d.db.Where("slug = ?", slug).First(&stats)
	// if res.Error != nil {
	// 	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
	// 		return nil, ErrNotFound
	// 	}
	// 	return nil, res.Error
	// }
	// return &stats, nil
}
