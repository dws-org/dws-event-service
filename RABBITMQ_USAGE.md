# RabbitMQ Integration Guide

This guide explains how to use RabbitMQ in your event service.

## Prerequisites

1. RabbitMQ running in Kubernetes (or locally)
2. Port forwarding set up for both AMQP and Management ports

## Port Forwarding

You need to forward both ports:

```bash
# AMQP port (for message publishing/consuming)
kubectl port-forward -n rabbitmq svc/rabbitmq-service 5672:5672

# Management UI port (for monitoring)
kubectl port-forward -n rabbitmq svc/rabbitmq-service-management 15672:15672
```

## Configuration

Update `configs/config.yaml` with your RabbitMQ settings:

```yaml
rabbitmq:
  enabled: true
  host: "localhost"  # Use 'localhost' when port-forwarding
  port: "5672"
  management_port: "15672"
  username: "guest"  # Default credentials
  password: "guest"
  virtual_host: "/"
```

## Usage Examples

### Publishing Events

```go
import (
    "github.com/oskargbc/dws-event-service.git/internal/pkg/rabbitmq"
    "github.com/google/uuid"
)

// Simple event publishing
eventData := map[string]interface{}{
    "event_id": "123",
    "name": "Concert",
    "location": "Stockholm",
}

err := rabbitmq.PublishEvent(
    "events",           // exchange name
    "event.created",    // routing key
    eventData,          // event data
)

// Using structured event messages
eventMsg := rabbitmq.NewEventMessage(
    "event.created",
    uuid.New().String(),
    "event-service",
    eventData,
)
eventMsg.Metadata["user_id"] = "user-123"

err := rabbitmq.PublishEventMessage("events", "event.created", eventMsg)
```

### Setting Up Exchanges and Queues

```go
import "github.com/oskargbc/dws-event-service.git/internal/pkg/rabbitmq"

// Set up exchange and queue with binding
err := rabbitmq.SetupExchangeAndQueue(
    "events",           // exchange name
    "topic",            // exchange type (direct, topic, fanout)
    "event-queue",      // queue name
    "event.*",          // routing key pattern
    true,               // durable
)
```

### Consuming Events

```go
import "github.com/oskargbc/dws-event-service.git/internal/pkg/rabbitmq"

// Define a message handler
handler := func(msg *rabbitmq.EventMessage) error {
    // Process the event
    fmt.Printf("Received event: %s\n", msg.EventType)
    fmt.Printf("Event data: %+v\n", msg.Data)
    
    // Your business logic here
    
    return nil // Return error to requeue message
}

// Start consuming
err := rabbitmq.StartEventConsumer(
    "event-queue",      // queue name
    "event-consumer",   // consumer tag
    handler,            // message handler
)
```

### Direct Service Usage

For more control, use the RabbitMQ service directly:

```go
import (
    "github.com/oskargbc/dws-event-service.git/internal/services"
    "github.com/rabbitmq/amqp091-go"
)

rabbitmqService := services.GetRabbitMQServiceInstance()
if rabbitmqService == nil {
    // RabbitMQ not enabled or failed to connect
    return
}

// Get channel for advanced operations
channel, err := rabbitmqService.GetChannel()
if err != nil {
    return err
}

// Declare exchange
err = rabbitmqService.DeclareExchange(
    "events",
    "topic",
    true,  // durable
    false, // autoDelete
    false, // internal
    false, // noWait
    nil,   // args
)

// Publish message
err = rabbitmqService.Publish(
    "events",
    "event.created",
    false, // mandatory
    false, // immediate
    amqp091.Publishing{
        ContentType:  "application/json",
        DeliveryMode: amqp091.Persistent,
        Body:        []byte(`{"event": "created"}`),
    },
)
```

## Health Checks

The health check endpoint (`/ready`) includes RabbitMQ status:

```bash
curl http://localhost:6906/ready
```

Response includes RabbitMQ status:
```json
{
  "status": "ok",
  "service": "event-service",
  "version": "0.1.0",
  "checks": {
    "database": {
      "status": "ok"
    },
    "rabbitmq": {
      "status": "ok"
    }
  }
}
```

## Common Exchange Types

- **direct**: Routes messages to queues based on exact routing key match
- **topic**: Routes messages using pattern matching on routing keys (e.g., `event.*`, `*.created`)
- **fanout**: Broadcasts messages to all bound queues (ignores routing key)
- **headers**: Routes based on message headers instead of routing key

## Best Practices

1. **Durable Queues**: Use `durable: true` for queues that should survive broker restarts
2. **Message Persistence**: Use `DeliveryMode: amqp091.Persistent` for important messages
3. **Manual Acknowledgments**: Always use manual ack (`autoAck: false`) for reliable processing
4. **Error Handling**: Return errors from handlers to requeue messages, or use Nack for DLQ
5. **Connection Management**: The service handles reconnection automatically, but monitor logs

## Troubleshooting

### Connection Issues

- Verify port forwarding is active: `kubectl port-forward -n rabbitmq svc/rabbitmq-service 5672:5672`
- Check credentials in `config.yaml`
- Verify RabbitMQ is running: `kubectl get pods -n rabbitmq`

### Message Not Received

- Check exchange and queue are declared
- Verify routing key matches binding
- Check consumer is running and not crashed
- Review RabbitMQ management UI at `http://localhost:15672`

### Management UI Access

Access the RabbitMQ management UI:
- URL: `http://localhost:15672`
- Default credentials: `guest` / `guest`

