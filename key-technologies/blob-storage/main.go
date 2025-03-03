package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/samber/lo"
)

const (
	useSSL     = false
	expiration = time.Minute * 15
)

type Config struct {
	Bucket    string `env:"BUCKET"`
	Endpoint  string `env:"MINIO_URL"`
	AccessKey string `env:"ACCESS_KEY_ID"`
	SecretKey string `env:"SECRET_ACCESS_KEY"`
}

var config Config

func getClient() (*minio.Client, error) {
	return minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: useSSL,
	})
}

func getSignedUploadURL(ctx context.Context, objectName string) (string, error) {
	minioClient, err := getClient()
	if err != nil {
		return "", err
	}

	presignedURL, err := minioClient.PresignedPutObject(ctx, config.Bucket, objectName, expiration)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

func uploadFile(url, data string) error {
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream") // Set appropriate content type

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func getSignedDownloadURL(ctx context.Context, objectName string) (string, error) {
	minioClient, err := getClient()
	if err != nil {
		return "", err
	}

	presignedURL, err := minioClient.PresignedGetObject(ctx, config.Bucket, objectName, expiration, nil)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

func downloadFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	ctx := context.Background()

	err := env.Parse(&config)
	if err != nil {
		log.Fatalf("Error parsing environment variables: %v", err)
	}

	objectName := fmt.Sprintf("file-%d.txt", time.Now().Unix())

	uploadURL := lo.Must(getSignedUploadURL(ctx, objectName))
	lo.Must0(uploadFile(uploadURL, "Hello, World!"+time.Now().String()))
	slog.Info("File uploaded successfully")

	downloadURL := lo.Must(getSignedDownloadURL(ctx, objectName))
	content := lo.Must(downloadFile(downloadURL))
	slog.Info("Downloaded content", "content", content)
}
