package db

import (
	"context"
	"errors"
	"log/slog"
	"math/rand"
	"sync/atomic"

	"github.com/ericbutera/system-design/labs/url-shortener/internal/base62"
	"github.com/ericbutera/system-design/labs/url-shortener/internal/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrDuplicate = errors.New("duplicate slug")
	ErrNotFound  = errors.New("not found")
	Counter      atomic.Uint64
)

type DB struct {
	db  *gorm.DB
	rdb *redis.Client
}

func New(dsn string, rdb *redis.Client) (*DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return &DB{
		db:  db,
		rdb: rdb,
	}, err
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

// Fetch the long url paired with a read-through cache
func (d *DB) GetURL(ctx context.Context, slug string, cache bool) (string, error) {
	slog.Debug("Geturl", "slug", slug)
	var url string

	if cache {
		url, err := d.rdb.Get(ctx, slug).Result()
		if err == nil && url != "" {
			slog.Debug("cache hit", "slug", slug, "url", url)
			return url, nil
		}
	}

	slog.Debug("cache miss", "slug", slug)
	res := d.db.Raw(`SELECT long FROM urls_v1 WHERE slug = ?`, slug).Scan(&url)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return "", ErrNotFound
		}
		return "", res.Error
	}

	if cache {
		d.rdb.Set(ctx, slug, url, 0)
		slog.Debug("caching", "slug", slug, "url", url)
	}

	slog.Debug("get url", "url", url)
	return url, nil
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
