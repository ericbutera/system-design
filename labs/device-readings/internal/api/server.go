package api

import (
	"log/slog"
	"net/http"

	"device-readings/internal/readings/queue"
	"device-readings/internal/readings/repo"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

type Server struct {
	router *gin.Engine
}

func New(writer queue.BatchReadingWriter, repo repo.Repo) (*Server, error) {
	handlers, err := NewHandlers(writer, repo)
	if err != nil {
		return nil, err
	}

	router := NewRouter(handlers)
	if err := router.SetTrustedProxies([]string{}); err != nil {
		return nil, err
	}

	return &Server{
		router: router,
	}, nil
}

func (s *Server) Start() error {
	return s.router.Run()
}

func NewRouter(handlers *Handlers) *gin.Engine {
	router := gin.New()
	routes(router, handlers)
	return router
}

func routes(router *gin.Engine, handlers *Handlers) {
	router.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	v1 := router.Group("/v1").
		Use(sloggin.New(slog.Default())).
		Use(gin.Recovery())
	{
		v1.POST("/readings", handlers.StoreReadings)
		v1.GET("/readings", handlers.GetReadings)
	}
}
