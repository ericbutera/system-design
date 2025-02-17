package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"device-readings/internal/readings/models"
	"device-readings/internal/readings/queue"
	"device-readings/internal/readings/repo"

	"github.com/gin-gonic/gin"
)

const DateFormatISO8601 = "2006-01-02T15:04:05-07:00"

type Handlers struct {
	producer queue.Producer
	repo     repo.Repo
}

func NewHandlers(producer queue.Producer, repo repo.Repo) (*Handlers, error) {
	return &Handlers{
		producer: producer,
		repo:     repo,
	}, nil
}

type ErrorResponse struct {
	Error string `json:"error"`
	// TODO: add error code
	// TODO: add request ID for correlating traces, logs, etc
}

type StoreReadingsRequest struct {
	Readings []models.BatchReading `binding:"required,dive" description:"Readings" json:"readings"`
}

type StoreReadingsResponse struct {
	ID string `description:"ID to check processing status" json:"id"`
}

// store readings for a device; will create a device if it doesn't exist.
func (h *Handlers) StoreReadings(c *gin.Context) {
	// TODO: validate path device matches body device
	// TODO: go validation library supports "friendly" error messages (i11n becomes a concern)
	// TODO: throw out batch if threshold of invalid readings are reached
	var req StoreReadingsRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: fmt.Sprintf("validation error: %s", err.Error()),
		})
		return
	}

	res, err := h.producer.Write(c.Request.Context(), req.Readings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "unable to save readings"}) // TODO: support validation errors; if errors.As(err, &repo.ValidationErrors{}) { c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid record: " + err.Error()}) return }
		return
	}

	c.JSON(http.StatusOK, &StoreReadingsResponse{
		ID: res.ID,
	})
}

type GetReadingsResponse struct {
	Readings JSONSlice[models.Reading] `json:"readings"`
}

func (h *Handlers) GetReadings(c *gin.Context) {
	res, err := h.repo.GetReadings(c.Request.Context(), repo.Filters{
		DeviceID: c.Query("device_id"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "unable to get readings"})
		return
	}
	c.JSON(http.StatusOK, &GetReadingsResponse{
		Readings: JSONSlice[models.Reading](res),
	})
}

func TimeFromString(s string) (time.Time, error) {
	return time.Parse(DateFormatISO8601, s)
}

func TimeToString(t time.Time) string {
	return t.Format(DateFormatISO8601)
}

// Prevents empty slices from being marshaled as `null` which can break some clients.
type JSONSlice[T any] []T

func (s JSONSlice[T]) MarshalJSON() ([]byte, error) {
	if s == nil {
		return json.Marshal([]T{})
	}
	return json.Marshal([]T(s))
}
