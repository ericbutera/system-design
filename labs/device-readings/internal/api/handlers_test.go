package api_test

import (
	"device-readings/internal/api"
	"device-readings/internal/readings/models"
	"device-readings/internal/readings/queue"
	"device-readings/internal/readings/repo"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

const (
	ReadingEndpoint = "/v1/readings"
)

type testSetup struct {
	handlers *api.Handlers
	router   *gin.Engine
	repo     *repo.MockRepo
}

func setup(t *testing.T) *testSetup {
	t.Helper()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	slog.SetDefault(logger)

	gin.SetMode(gin.TestMode)

	producer := new(queue.MockProducer)
	repo := new(repo.MockRepo)

	handlers, err := api.NewHandlers(producer, repo)
	require.NoError(t, err)

	router := api.NewRouter(handlers)

	return &testSetup{
		handlers: handlers,
		router:   router,
		// repo:       repo,
	}
}

func TestStoreReadings(t *testing.T) {
	s := setup(t)

	raw := `
	{
		"readings": [
			{"device_id":"36d5658a-6908-479e-887e-a949ec199272", "type": "temperature", "timestamp":"2021-09-01T17:00:00-05:00","value":17.3}
		]
	}
	`
	ts, err := api.TimeFromString("2021-09-01T17:00:00-05:00")
	require.NoError(t, err)

	s.repo.EXPECT().
		StoreReadings([]*models.Reading{
			{DeviceID: "36d5658a-6908-479e-887e-a949ec199272", ReadingType: "temperature", Timestamp: ts, Value: 17.3},
		}).
		Return(repo.StoreReadingsResult{}, nil)

	request := httptest.NewRequest(http.MethodPost, ReadingEndpoint, strings.NewReader(raw))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	s.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusNoContent, recorder.Code)
	t.Skip("TODO: validate response using unmarshal")
}

// TODO TestGetReadings

// func unmarshal(t *testing.T, recorder *httptest.ResponseRecorder, v any) {
// 	t.Helper()
// 	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), v))
// }
