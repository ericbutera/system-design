package api

import (
	"errors"
	"fmt"
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
	ApiUrl    string `env:"API_URL" envDefault:"http://api:8080"`
	DSN       string `env:"DSN"`
	RedisAddr string `env:"REDIS_ADDR" envDefault:"redis-master:6379"`
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
	router := gin.New() // router.Use(gin.Recovery())

	router.POST("/noop", func(c *gin.Context) { c.Status(http.StatusOK) })
	router.POST("/v1/urls", a.CreateURL_PG_Counter)
	router.POST("/beta/urls/random-generator", a.CreateURL_RandomGenerator)
	router.POST("/beta/urls/in-proc", a.CreateURL_AtomicCounter)

	router.GET("/v1/:slug", a.Redirect)

	// TODO: router.GET("/v1/urls/:slug/stats", a.GetStats)

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

func URLFromRequest(c *gin.Context) (*models.URL, error) {
	// TODO: sanitize url (normalize for deduplication, prevent XSS, open redirect)
	var req CreateURLRequest
	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		return nil, err
	}
	return &models.URL{
		Long:      req.Long,
		Slug:      req.Slug,
		ExpiresAt: req.ExpiresAt,
		UserID:    UserID,
	}, nil
}

func handleCreate(c *gin.Context, createFn func(*models.URL) error, respFn func(*models.URL) *CreateURLResponse) {
	url, err := URLFromRequest(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	if err := createFn(url); err != nil {
		if errors.Is(err, db.ErrDuplicate) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "url already exists"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to create url"})
		return
	}
	c.JSON(http.StatusOK, respFn(url))
}

func (a *API) handleCreateResponse(url *models.URL) *CreateURLResponse {
	return &CreateURLResponse{
		Slug:      url.Slug,
		Short:     fmt.Sprintf("%s/%s", a.config.ApiUrl, url.Slug),
		Long:      url.Long,
		ExpiresAt: url.ExpiresAt,
	}
}

func (a *API) CreateURL_PG_Counter(c *gin.Context) {
	handleCreate(c, a.repo.CreateURL_PG_Counter, a.handleCreateResponse)
}

func (a *API) CreateURL_RandomGenerator(c *gin.Context) {
	handleCreate(c, a.repo.CreateURL_RandomGenerator, a.handleCreateResponse)
}

func (a *API) CreateURL_AtomicCounter(c *gin.Context) {
	handleCreate(c, a.repo.CreateURL_AtomicCounter, a.handleCreateResponse)
}

func (a *API) Redirect(c *gin.Context) {
	slug := c.Param("slug")
	cache := !c.Request.URL.Query().Has("nocache")
	url, err := a.repo.GetURL(c.Request.Context(), slug, cache)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "url not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to get url"})
		return
	}

	// TODO: increment stats

	c.Redirect(http.StatusFound, url)
}
