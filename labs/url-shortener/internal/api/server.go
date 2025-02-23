package api

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/ericbutera/system-design/labs/url-shortener/internal/db"
	"github.com/ericbutera/system-design/labs/url-shortener/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

const (
	UserID = 1
)

type API struct {
	config Config
	repo   *db.DB
}

func New(config Config, repo *db.DB) *API {
	return &API{
		config: config,
		repo:   repo,
	}
}

type Config struct {
	ApiUrl string `env:"API_URL" envDefault:"http://api:8080"`
	DSN    string `env:"DSN"`
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

func (a *API) Server() *gin.Engine {
	router := gin.Default()

	router.POST("/v1/urls", a.CreateURL)
	router.GET("/:slug", a.Redirect)
	router.GET("/v1/urls/:slug/stats", a.GetStats)

	return router
}

func (a *API) GetStats(c *gin.Context) {
	stats, err := a.repo.GetStats(c.Param("slug"))
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
}

func (a *API) CreateURL(c *gin.Context) {
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
	err := a.repo.CreateURL(url)
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
		Short: a.config.ApiUrl + "/" + url.Slug,
		Long:  url.Long,
	})
}

func (a *API) Redirect(c *gin.Context) {
	url, err := a.repo.GetURL(c.Param("slug"))
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "url not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to get url"})
		return
	}
	slog.Debug("redirecting", "slug", url.Slug, "long", url.Long)

	// TODO: increment stats
	c.Redirect(http.StatusFound, url.Long)
}
