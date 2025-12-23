package rabbitmq

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oskargbc/dws-event-service.git/configs"
	"github.com/oskargbc/dws-event-service.git/internal/pkg/rabbitmq"
	"github.com/oskargbc/dws-event-service.git/internal/services"
)

// Controller handles RabbitMQ test endpoints
type Controller struct {
	rabbitmqService *services.RabbitMQService
	config          *configs.Config
}

// NewController creates a new RabbitMQ test controller
func NewController() *Controller {
	envConfig := configs.GetEnvConfig()
	var rabbitmqService *services.RabbitMQService
	if envConfig.RabbitMQ.Enabled {
		rabbitmqService = services.GetRabbitMQServiceInstance()
	}

	return &Controller{
		rabbitmqService: rabbitmqService,
		config:          envConfig,
	}
}

// TestConnection tests the RabbitMQ connection
// @Summary      Test RabbitMQ connection
// @Description  Tests if RabbitMQ connection is working
// @Tags         rabbitmq
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      503  {object}  map[string]interface{}
// @Router       /rabbitmq/test [get]
func (rc *Controller) TestConnection(c *gin.Context) {
	if rc.rabbitmqService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "RabbitMQ service is not available or not enabled",
		})
		return
	}

	ctx := c.Request.Context()
	if err := rc.rabbitmqService.HealthCheck(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "RabbitMQ health check failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "RabbitMQ connection is working",
		"config": gin.H{
			"host":     rc.config.RabbitMQ.Host,
			"port":     rc.config.RabbitMQ.Port,
			"enabled":  rc.config.RabbitMQ.Enabled,
			"username": rc.config.RabbitMQ.Username,
		},
	})
}

// PublishTestMessageRequest represents the request to publish a test message
type PublishTestMessageRequest struct {
	Exchange   string                 `json:"exchange" binding:"required"`
	RoutingKey string                 `json:"routingKey" binding:"required"`
	Message    map[string]interface{} `json:"message"`
}

// PublishTestMessage publishes a test message to RabbitMQ
// @Summary      Publish a test message
// @Description  Publishes a test message to RabbitMQ exchange
// @Tags         rabbitmq
// @Accept       json
// @Produce      json
// @Param        request  body      PublishTestMessageRequest  true  "Publish request"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      503      {object}  map[string]interface{}
// @Router       /rabbitmq/publish [post]
func (rc *Controller) PublishTestMessage(c *gin.Context) {
	if rc.rabbitmqService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "RabbitMQ service is not available or not enabled",
		})
		return
	}

	var req PublishTestMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request payload",
			"error":   err.Error(),
		})
		return
	}

	// Create a structured event message
	eventMsg := rabbitmq.NewEventMessage(
		"test.message",
		uuid.New().String(),
		"event-service",
		req.Message,
	)
	eventMsg.Metadata["timestamp"] = time.Now().UTC().Format(time.RFC3339)
	eventMsg.Metadata["test"] = true

	// Publish the message
	if err := rabbitmq.PublishEventMessage(req.Exchange, req.RoutingKey, eventMsg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to publish message",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       "ok",
		"message":      "Message published successfully",
		"event_id":     eventMsg.EventID,
		"exchange":     req.Exchange,
		"routing_key":  req.RoutingKey,
		"published_at": eventMsg.Timestamp,
	})
}

// SetupTestExchangeAndQueue sets up a test exchange and queue
// @Summary      Setup test exchange and queue
// @Description  Creates a test exchange and queue with binding
// @Tags         rabbitmq
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      503  {object}  map[string]interface{}
// @Router       /rabbitmq/setup [post]
func (rc *Controller) SetupTestExchangeAndQueue(c *gin.Context) {
	if rc.rabbitmqService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "RabbitMQ service is not available or not enabled",
		})
		return
	}

	// Setup a test exchange and queue
	exchangeName := "test-events"
	queueName := "test-queue"
	routingKey := "test.*"

	if err := rabbitmq.SetupExchangeAndQueue(
		exchangeName,
		"topic",
		queueName,
		routingKey,
		true, // durable
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to setup exchange and queue",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "ok",
		"message":     "Exchange and queue setup successfully",
		"exchange":    exchangeName,
		"queue":       queueName,
		"routing_key": routingKey,
	})
}
