package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
)

type Event struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

func main() {
	r := gin.Default()

	// No need for cors if using ingress
	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"*"}, // Allow all origins
	// 	AllowMethods:     []string{"GET", "POST", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// }))

	r.GET("/events", func(c *gin.Context) {
		header := c.Writer.Header()
		header.Set("Content-Type", "text/event-stream")
		header.Set("Cache-Control", "no-cache")
		header.Set("Connection", "keep-alive")
		header.Set("Transfer-Encoding", "chunked")

		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming unsupported"})
			return
		}

		counter := 0
		for range time.Tick(1 * time.Second) {
			event := Event{
				Timestamp: time.Now().Format(time.RFC3339),
				Message:   faker.Sentence(),
			}

			data, err := json.Marshal(event)
			if err != nil {
				slog.Error("Failed to marshal JSON", "error", err)
				continue
			}

			var msg string
			if counter%2 == 0 {
				msg = "event: important\ndata: " + string(data) + "\n\n"
			} else {
				msg = "data: " + string(data) + "\n\n"
			}

			_, err = c.Writer.Write([]byte(msg))
			if err != nil {
				slog.Error("Failed to write data", "error", err)
				continue
			}

			flusher.Flush()
			slog.Info("Event sent", "event", event)
			counter++
		}
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
