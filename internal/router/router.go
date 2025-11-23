package router

import (
	"github.com/oskargbc/dws-event-service.git/internal/controllers/events"
	"github.com/oskargbc/dws-event-service.git/internal/controllers/health"
	"github.com/oskargbc/dws-event-service.git/internal/middlewares"
	"github.com/oskargbc/dws-event-service.git/internal/pkg/logger"

	"github.com/gin-gonic/gin"
)

func NewGinRouter(mode string) *gin.Engine {

	gin.SetMode(mode)

	gin.DefaultWriter = logger.NewLogrusLogger().Writer()

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, x-api-key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.Use(middlewares.ErrorHandle())

	healthController := health.NewController()
	router.GET("/livez", healthController.Live)
	router.GET("/readyz", healthController.Ready)
	router.GET("/healthz", healthController.Ready)
	router.GET("/_meta", healthController.Info)

	//	router.Use(APIKeyAuthMiddleware())

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		eventsController := events.NewController()
		v1.GET("/events", eventsController.GetEvents)
		v1.GET("/events/:id", eventsController.GetEventByID)
	}

	return router
}
