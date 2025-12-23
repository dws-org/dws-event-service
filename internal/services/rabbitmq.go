package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/oskargbc/dws-event-service.git/configs"
	"github.com/oskargbc/dws-event-service.git/internal/pkg/logger"

	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

var rabbitmqInstance *RabbitMQService
var rabbitmqLock = &sync.Mutex{}

type RabbitMQService struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	logger  *logrus.Logger
	config  *configs.RabbitMQ
}

// GetRabbitMQServiceInstance returns a singleton instance of RabbitMQService
func GetRabbitMQServiceInstance() *RabbitMQService {
	if rabbitmqInstance == nil {
		rabbitmqLock.Lock()
		defer rabbitmqLock.Unlock()
		if rabbitmqInstance == nil {
			envConfig := configs.GetEnvConfig()
			if !envConfig.RabbitMQ.Enabled {
				return nil
			}

			rabbitmqInstance = &RabbitMQService{
				logger: logger.NewLogrusLogger(),
				config: &envConfig.RabbitMQ,
			}
			if err := rabbitmqInstance.connect(); err != nil {
				rabbitmqInstance.logger.Fatalf("Failed to connect to RabbitMQ: %v", err)
			}
		}
	}

	return rabbitmqInstance
}

// connect establishes connection to RabbitMQ
func (r *RabbitMQService) connect() error {
	if r.config == nil {
		return errors.New("RabbitMQ config is not initialized")
	}

	// Build connection URL
	url := fmt.Sprintf("amqp://%s:%s@%s:%s%s",
		r.config.Username,
		r.config.Password,
		r.config.Host,
		r.config.Port,
		r.config.VirtualHost,
	)

	r.logger.Infof("Connecting to RabbitMQ at %s:%s", r.config.Host, r.config.Port)

	var err error
	r.conn, err = amqp091.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	r.channel, err = r.conn.Channel()
	if err != nil {
		r.conn.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	r.logger.Infoln("Successfully connected to RabbitMQ")
	return nil
}

// EnsureConnected ensures the RabbitMQ connection is active, reconnecting if necessary
func (r *RabbitMQService) EnsureConnected(ctx context.Context) error {
	if r.conn == nil || r.conn.IsClosed() {
		r.logger.Warnln("RabbitMQ connection is closed, reconnecting...")
		return r.connect()
	}

	if r.channel == nil || r.channel.IsClosed() {
		r.logger.Warnln("RabbitMQ channel is closed, recreating...")
		var err error
		r.channel, err = r.conn.Channel()
		if err != nil {
			return fmt.Errorf("failed to recreate channel: %w", err)
		}
	}

	return nil
}

// GetChannel returns the RabbitMQ channel, ensuring connection is active
func (r *RabbitMQService) GetChannel() (*amqp091.Channel, error) {
	ctx := context.Background()
	if err := r.EnsureConnected(ctx); err != nil {
		return nil, err
	}
	return r.channel, nil
}

// Publish publishes a message to an exchange
func (r *RabbitMQService) Publish(exchange, routingKey string, mandatory, immediate bool, msg amqp091.Publishing) error {
	channel, err := r.GetChannel()
	if err != nil {
		return err
	}

	return channel.Publish(exchange, routingKey, mandatory, immediate, msg)
}

// DeclareQueue declares a queue if it doesn't exist
func (r *RabbitMQService) DeclareQueue(name string, durable, autoDelete, exclusive, noWait bool, args amqp091.Table) (amqp091.Queue, error) {
	channel, err := r.GetChannel()
	if err != nil {
		return amqp091.Queue{}, err
	}

	return channel.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
}

// DeclareExchange declares an exchange if it doesn't exist
func (r *RabbitMQService) DeclareExchange(name, kind string, durable, autoDelete, internal, noWait bool, args amqp091.Table) error {
	channel, err := r.GetChannel()
	if err != nil {
		return err
	}

	return channel.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
}

// QueueBind binds a queue to an exchange
func (r *RabbitMQService) QueueBind(queue, key, exchange string, noWait bool, args amqp091.Table) error {
	channel, err := r.GetChannel()
	if err != nil {
		return err
	}

	return channel.QueueBind(queue, key, exchange, noWait, args)
}

// Consume starts consuming messages from a queue
func (r *RabbitMQService) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp091.Table) (<-chan amqp091.Delivery, error) {
	channel, err := r.GetChannel()
	if err != nil {
		return nil, err
	}

	return channel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}

// HealthCheck validates that RabbitMQ connection is active
func (r *RabbitMQService) HealthCheck(ctx context.Context) error {
	if r.conn == nil {
		return errors.New("RabbitMQ connection is not initialized")
	}

	if r.conn.IsClosed() {
		return errors.New("RabbitMQ connection is closed")
	}

	// Try to ensure connection is active
	if err := r.EnsureConnected(ctx); err != nil {
		return fmt.Errorf("RabbitMQ health check failed: %w", err)
	}

	// Try to declare a temporary queue to verify channel is working
	channel, err := r.GetChannel()
	if err != nil {
		return fmt.Errorf("RabbitMQ channel health check failed: %w", err)
	}

	// Declare a temporary queue to verify everything works
	testQueue, err := channel.QueueDeclare(
		"",    // name (empty = auto-generate)
		false, // durable
		true,  // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("RabbitMQ queue declaration test failed: %w", err)
	}

	// Delete the test queue
	_, err = channel.QueueDelete(testQueue.Name, false, false, false)
	if err != nil {
		r.logger.Warnf("Failed to delete test queue: %v", err)
	}

	return nil
}

// Close closes the RabbitMQ connection and channel
func (r *RabbitMQService) Close() error {
	var errs []error

	if r.channel != nil && !r.channel.IsClosed() {
		if err := r.channel.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close channel: %w", err))
		}
	}

	if r.conn != nil && !r.conn.IsClosed() {
		if err := r.conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close connection: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing RabbitMQ: %v", errs)
	}

	r.logger.Infoln("RabbitMQ connection closed")
	return nil
}

// StartConsumer starts a consumer that processes messages from a queue
func (r *RabbitMQService) StartConsumer(queueName, consumerTag string, handler func(amqp091.Delivery)) error {
	channel, err := r.GetChannel()
	if err != nil {
		return err
	}

	msgs, err := channel.Consume(
		queueName,   // queue
		consumerTag, // consumer tag
		false,       // auto-ack (set to false to manually acknowledge)
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		r.logger.Infof("Starting consumer for queue: %s", queueName)
		for msg := range msgs {
			handler(msg)
		}
		r.logger.Warnf("Consumer for queue %s stopped", queueName)
	}()

	return nil
}

// PublishJSON publishes a JSON message to an exchange
func (r *RabbitMQService) PublishJSON(exchange, routingKey string, body interface{}) error {
	channel, err := r.GetChannel()
	if err != nil {
		return err
	}

	// Convert body to JSON bytes
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	msg := amqp091.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp091.Persistent, // Make message persistent
		Timestamp:    time.Now(),
		Body:         jsonBody,
	}

	return channel.Publish(exchange, routingKey, false, false, msg)
}
