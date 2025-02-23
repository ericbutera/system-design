package db

import (
	"errors"
	"math/rand"
	"sync/atomic"

	"github.com/ericbutera/system-design/labs/url-shortener/internal/base62"
	"github.com/ericbutera/system-design/labs/url-shortener/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrDuplicate = errors.New("duplicate slug")
	ErrNotFound  = errors.New("not found")
	Counter      atomic.Uint64
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

func (d *DB) GetURL(slug string) (*models.URL, error) {
	var url_v1 models.URL
	res := d.db.Table("urls_v1").Where("slug = ?", slug).First(&url_v1)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, res.Error
	}
	return &url_v1, nil
}

func (d *DB) GetStats(slug string) (*models.URLStats, error) {
	panic("not implemented")
	// var stats models.URLStats
	// res := d.db.Table("urls_v1").Where("slug = ?", slug).First(&stats)
	// if res.Error != nil {
	// 	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
	// 		return nil, ErrNotFound
	// 	}
	// 	return nil, res.Error
	// }
	// return &stats, nil
}

// Random generator. Collision risk.
func (d *DB) CreateURL_RandomGenerator(url *models.URL) error {
	counter := rand.Intn(1_000_000_000_000)
	slug := base62.EncodeInt64(int64(counter))
	url.Slug = slug
	return d.createURL("urls_v0", url)
}

// Counter using PG sequence. Durable and safe. Unsure of performance.
func (d *DB) CreateURL_PG_Counter(url *models.URL) error {
	if url.Slug == "" {
		slug, err := d.GenerateSlug()
		if err != nil {
			return err
		}
		url.Slug = slug.Slug
	}
	return d.createURL("urls_v1", url)
}

// In-process counter. Everything breaks if the app crashes or needs autoscaling.
func (d *DB) CreateURL_AtomicCounter(url *models.URL) error {
	Counter.Add(1)
	counter := Counter.Load()
	slug := base62.EncodeInt64(int64(counter))
	url.Slug = slug
	return d.createURL("urls_v2", url)
}

func (d *DB) createURL(table string, url *models.URL) error {
	// TODO: support expires_at
	sql := "INSERT INTO " + table + " (slug, long, user_id, expires_at) VALUES (?, ?, ?, ?)"
	res := d.db.Exec(sql, url.Slug, url.Long, url.UserID, url.ExpiresAt)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func init() {
	Counter.Add(1)
}
