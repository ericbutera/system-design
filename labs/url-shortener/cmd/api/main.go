package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/ericbutera/system-design/labs/url-shortener/internal/db"
	"github.com/ericbutera/system-design/labs/url-shortener/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

const (
	UserID = 1
	ApiUrl = "http://api:8080"
)

type Config struct {
	DSN string `env:"DSN"`
}

type CreateURLRequest struct {
	models.URL
}

type CreateURLResponse struct {
	Slug      string     `json:"slug"`
	Short     string     `json:"short"`
	Long      string     `json:"long"`
	ExpiresAt *time.Time `json:"expires_at"`
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var config Config
	lo.Must0(env.Parse(&config))

	repo := lo.Must(db.New(config.DSN))
	server := Server(repo)

	srvErr := make(chan error, 1)
	go func() { srvErr <- server.Run() }()

	select {
	case err := <-srvErr:
		slog.Error("server error", "error", err)
		os.Exit(1)
	case <-ctx.Done():
		slog.Info("shutting down")
	}
}

func Server(repo *db.DB) *gin.Engine {
	router := gin.Default()

	router.POST("/v1/urls", func(c *gin.Context) {
		// TODO: sanitize url (normalize for deduplication, prevent XSS, open redirect)
		var req CreateURLRequest
		if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
			return
		}

		url := &models.URL{
			Long:      req.Long,
			Slug:      req.Slug,
			ExpiresAt: req.ExpiresAt,
			UserID:    UserID,
		}

		if url.Slug == "" {
			slug, err := repo.GenerateSlug()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to generate short url"})
				return
			}
			slog.Info("generated slug", "slug", slug.Slug, "counter", slug.Counter)
			url.Slug = slug.Slug
		}

		err := repo.CreateURL(url)
		if err != nil {
			if errors.Is(err, db.ErrDuplicate) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "url already exists"})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to create url"})
			return
		}
		c.JSON(http.StatusOK, &CreateURLResponse{
			Slug:  url.Slug,
			Short: ApiUrl + "/" + url.Slug,
			Long:  url.Long,
		})
	})

	router.GET("/:slug", func(c *gin.Context) {
		url, err := repo.GetURL(c.Param("slug"))
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "url not found"})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to get url"})
			return
		}
		slog.Info("redirecting", "slug", url.Slug, "long", url.Long)

		// TODO: increment stats
		c.Redirect(http.StatusFound, url.Long)
	})

	router.GET("/v1/urls/:slug/stats", func(c *gin.Context) {
		stats, err := repo.GetStats(c.Param("slug"))
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "url not found"})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to get stats"})
			return
		}

		// TODO: check if user is authorized to view stats
		// if stats.UserID != UserID {
		// 	c.AbortWithStatus(http.StatusForbidden)
		// 	return
		// }

		c.JSON(http.StatusOK, stats)
	})
	return router
}
