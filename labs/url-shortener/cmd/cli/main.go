package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/sourcegraph/conc/pool"
	"resty.dev/v3"
)

const (
	Workers       = 100
	RecordCount   = 50
	ApiUrl        = "http://api:8080"
	FakeUrl       = `http://localhost?t=%d_%d`
	CreateRequest = `{"long": "%s"}`
)

var (
	ErrInternal = errors.New("internal error")
)

func main() {
	benchCreate(1_000)
	benchRead(1_000)
}

func benchCreate(count int) {
	client := resty.New()
	defer client.Close()

	bench := []struct {
		name     string
		endpoint string
	}{
		{"create noop", "/noop"},
		{"create v0 - Random number generator", "/v0/urls"},
		{"create v1 - PG counter", "/v1/urls"},
		{"create v2 - Atomic counter", "/v2/urls"},
		{"create v3 - Redis counter", "/v3/urls"},
	}

	for _, b := range bench {
		slog.Info("benchmark", "name", b.name, "count", RecordCount)
		success := atomic.Int64{}
		errors := atomic.Int64{}
		start := time.Now()
		t := start.UnixMilli()

		p := pool.New().WithMaxGoroutines(Workers)
		for x := 0; x < count; x++ {
			p.Go(func() {
				long := fmt.Sprintf(FakeUrl, t, x)
				_, err := CreateURL(client, b.endpoint, long)
				if err != nil {
					errors.Add(1)
				} else {
					success.Add(1)
				}
			})
		}
		p.Wait()

		end := time.Now()
		slog.Info("bench",
			"duration", end.Sub(start),
			"success", success.Load(),
			"errors", errors.Load(),
		)
	}
}

type URL struct {
	Short     string     `json:"short"`
	Long      string     `json:"long"`
	Slug      string     `json:"slug"`
	ExpiresAt *time.Time `json:"expires_at"`
}

func CreateURL(client *resty.Client, endpoint string, long string) (string, error) {
	var url URL
	res, err := client.R().
		SetBody(fmt.Sprintf(CreateRequest, long)).
		SetResult(&url).
		Post(ApiUrl + endpoint)

	if err != nil {
		return "", err
	}

	code := res.StatusCode()
	if code != http.StatusOK && code != http.StatusCreated {
		return "", ErrInternal
	}

	return url.Short, nil
}

func benchRead(count int) {
	client := resty.New().SetRedirectPolicy(resty.NoRedirectPolicy())
	defer client.Close()

	slog.Info("benchmark", "name", "/v1/:slug", "count", count)
	success := atomic.Int64{}
	errors := atomic.Int64{}
	start := time.Now()

	p := pool.New().WithMaxGoroutines(Workers)
	for x := 0; x < count; x++ {
		p.Go(func() {
			_, err := GetUrl(client, "/v1/test")
			if err != nil {
				errors.Add(1)
			} else {
				success.Add(1)
			}
		})
	}
	p.Wait()

	end := time.Now()
	slog.Info("bench",
		"duration", end.Sub(start),
		"success", success.Load(),
		"errors", errors.Load(),
	)
}

func GetUrl(client *resty.Client, endpoint string) (string, error) {
	var url URL
	res, err := client.R().
		SetResult(&url).
		Get(ApiUrl + endpoint)

	if err != nil {
		slog.Info("error", "err", err)
		return "", err
	}

	code := res.StatusCode()
	if code != http.StatusFound {
		return "", ErrInternal
	}

	return url.Short, nil
}
