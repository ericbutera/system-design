package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/samber/lo"
	"resty.dev/v3"
)

const BatchSize = 1_000
const ApiUrl = "http://api:8080"
const CreateURI = "/v1/urls"
const CreateRequest = `{"long": "%s"}`

func main() {
	client := resty.New()
	defer client.Close()

	for x := 0; x < BatchSize; x++ {
		long := fmt.Sprintf(`http://localhost?t=%d_%d`, time.Now().UnixMilli(), x)
		short := lo.Must(CreateURL(client, long))
		slog.Info("url", "long", long, "short", short)
	}
}

type URL struct {
	Short     string     `json:"short"`
	Long      string     `json:"long"`
	Slug      string     `json:"slug"`
	ExpiresAt *time.Time `json:"expires_at"`
}

func CreateURL(client *resty.Client, long string) (string, error) {
	var url URL
	res, err := client.R().
		SetBody(fmt.Sprintf(CreateRequest, long)).
		SetResult(&url).
		Post(ApiUrl + CreateURI)

	if err != nil {
		return "", err
	}
	slog.Info("response", "code", res.StatusCode())

	return url.Short, nil
}
