package health

import (
	"github.com/oskargbc/dws-event-service.git/configs"
	"github.com/oskargbc/dws-event-service.git/internal/services"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Controller exposes liveness/readiness endpoints that can be reused across
// microservices.
type Controller struct {
	service   configs.Service
	startedAt time.Time
}

// NewController instantiates a health controller backed by the shared config.
func NewController() *Controller {
	envCfg := configs.GetEnvConfig()
	return &Controller{
		service:   envCfg.Service,
		startedAt: time.Now().UTC(),
	}
}

// Live responds with a basic liveness payload that can be used by container
// orchestrators.
func (hc *Controller) Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"service":   hc.service.Slug,
		"timestamp": time.Now().UTC(),
		"uptime":    time.Since(hc.startedAt).Round(time.Second).String(),
	})
}

// Ready validates downstream dependencies (database, message brokers, etc) and
// reports readiness status.
func (hc *Controller) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	dbStatus := "skipped"
	if dbService := services.GetDatabaseSeviceInstance(); dbService != nil {
		if err := dbService.HealthCheck(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "fail",
				"service": hc.service.Slug,
				"version": hc.service.Version,
				"checks": gin.H{
					"database": gin.H{
						"status":  "fail",
						"details": err.Error(),
					},
				},
			})
			return
		}
		dbStatus = "ok"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": hc.service.Slug,
		"version": hc.service.Version,
		"checks": gin.H{
			"database": gin.H{
				"status": dbStatus,
			},
		},
		"timestamp": time.Now().UTC(),
	})
}

// Info exposes service metadata for debugging and observability tooling.
func (hc *Controller) Info(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":        hc.service.Name,
		"slug":        hc.service.Slug,
		"description": hc.service.Description,
		"version":     hc.service.Version,
		"tags":        hc.service.Tags,
		"startedAt":   hc.startedAt,
		"uptime":      time.Since(hc.startedAt).Round(time.Second).String(),
	})
}
