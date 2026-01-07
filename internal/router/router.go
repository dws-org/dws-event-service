package router

import (
	"github.com/oskargbc/dws-event-service.git/docs"
	"github.com/oskargbc/dws-event-service.git/internal/controllers/events"
	"github.com/oskargbc/dws-event-service.git/internal/controllers/health"
	rabbitmqController "github.com/oskargbc/dws-event-service.git/internal/controllers/rabbitmq"
	"github.com/oskargbc/dws-event-service.git/internal/middlewares"
	"github.com/oskargbc/dws-event-service.git/internal/pkg/logger"
	"github.com/oskargbc/dws-event-service.git/internal/pkg/metrics"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewGinRouter(mode string) *gin.Engine {

	gin.SetMode(mode)

	// Configure Swagger base path to match our API prefix
	docs.SwaggerInfo.BasePath = "/api/v1"

	gin.DefaultWriter = logger.NewLogrusLogger().Writer()

	router := gin.Default()

	// Add Prometheus middleware
	router.Use(metrics.PrometheusMiddleware("dws-event-service"))

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

	// Prometheus metrics endpoint (no auth required)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	healthController := health.NewController()
	router.GET("/livez", healthController.Live)
	router.GET("/readyz", healthController.Ready)
	router.GET("/healthz", healthController.Ready)
	router.GET("/_meta", healthController.Info)

	/*
	* ONLY FOR TESTING PURPOSESr 
	*/

	// RabbitMQ test endpoints (no auth required for testing)
	rabbitmqTestController := rabbitmqController.NewController()
	router.GET("/rabbitmq/test", rabbitmqTestController.TestConnection)
	router.POST("/rabbitmq/publish", rabbitmqTestController.PublishTestMessage)
	router.POST("/rabbitmq/setup", rabbitmqTestController.SetupTestExchangeAndQueue)

	// Swagger UI and OpenAPI JSON endpoints (no auth)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))


	/*
	* LIVE
	*/
	//	router.Use(APIKeyAuthMiddleware())

	// API v1 routes (protected by Keycloak auth middleware)
	v1 := router.Group("/api/v1", middlewares.KeycloakAuthMiddleware())
	{
		eventsController := events.NewController()
		v1.GET("/events", eventsController.GetEvents)
		v1.GET("/events/:id", eventsController.GetEventByID)
		// Only users with the "Organiser" realm role may create events
		v1.POST("/events", middlewares.RequireRole("Organiser"), eventsController.CreateEvent)
	}

	return router
}
