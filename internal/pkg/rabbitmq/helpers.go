package rabbitmq

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/oskargbc/dws-event-service.git/internal/services"
	"github.com/rabbitmq/amqp091-go"
)

// PublishEvent publishes an event to RabbitMQ
func PublishEvent(exchange, routingKey string, event interface{}) error {
	rabbitmqService := services.GetRabbitMQServiceInstance()
	if rabbitmqService == nil {
		return fmt.Errorf("RabbitMQ service is not available")
	}

	return rabbitmqService.PublishJSON(exchange, routingKey, event)
}

// SetupExchangeAndQueue sets up an exchange and queue with binding
func SetupExchangeAndQueue(exchangeName, exchangeType, queueName, routingKey string, durable bool) error {
	rabbitmqService := services.GetRabbitMQServiceInstance()
	if rabbitmqService == nil {
		return fmt.Errorf("RabbitMQ service is not available")
	}

	// Declare exchange
	if err := rabbitmqService.DeclareExchange(
		exchangeName,
		exchangeType,
		durable,
		false, // autoDelete
		false, // internal
		false, // noWait
		nil,   // args
	); err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declare queue
	_, err := rabbitmqService.DeclareQueue(
		queueName,
		durable,
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind queue to exchange
	if err := rabbitmqService.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false, // noWait
		nil,   // args
	); err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	return nil
}

// EventMessage represents a standard event message structure
type EventMessage struct {
	EventType string                 `json:"event_type"`
	EventID   string                 `json:"event_id"`
	Timestamp time.Time              `json:"timestamp"`
	Source    string                 `json:"source"`
	Data      interface{}            `json:"data"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// NewEventMessage creates a new event message
func NewEventMessage(eventType, eventID, source string, data interface{}) *EventMessage {
	return &EventMessage{
		EventType: eventType,
		EventID:   eventID,
		Timestamp: time.Now().UTC(),
		Source:    source,
		Data:      data,
		Metadata:  make(map[string]interface{}),
	}
}

// PublishEventMessage publishes a structured event message
func PublishEventMessage(exchange, routingKey string, msg *EventMessage) error {
	return PublishEvent(exchange, routingKey, msg)
}

// MessageHandler is a function type for handling consumed messages
type MessageHandler func(msg *EventMessage) error

// StartEventConsumer starts consuming events from a queue
func StartEventConsumer(queueName, consumerTag string, handler MessageHandler) error {
	rabbitmqService := services.GetRabbitMQServiceInstance()
	if rabbitmqService == nil {
		return fmt.Errorf("RabbitMQ service is not available")
	}

	return rabbitmqService.StartConsumer(queueName, consumerTag, func(delivery amqp091.Delivery) {
		var eventMsg EventMessage
		if err := json.Unmarshal(delivery.Body, &eventMsg); err != nil {
			// Log error but don't acknowledge message so it can be retried or sent to DLQ
			fmt.Printf("Failed to unmarshal message: %v\n", err)
			delivery.Nack(false, false) // Don't requeue if it's a format error
			return
		}

		// Process the message
		if err := handler(&eventMsg); err != nil {
			fmt.Printf("Failed to process message: %v\n", err)
			delivery.Nack(false, true) // Requeue on processing error
			return
		}

		// Acknowledge successful processing
		delivery.Ack(false)
	})
}
