# RabbitMQ Testing Guide

This guide shows you how to test if RabbitMQ is working correctly in your event service.

## Prerequisites

1. **Port-forward RabbitMQ** (in a separate terminal):
   ```bash
   # AMQP port for messaging
   kubectl port-forward -n rabbitmq svc/rabbitmq-service 5672:5672
   
   # Management UI (optional, for monitoring)
   kubectl port-forward -n rabbitmq svc/rabbitmq-service-management 15672:15672
   ```

2. **Start your service**:
   ```bash
   export DATABASE_URL="postgresql://postgres:w5f7eE-CJrTjPl-DfYZkV@localhost:5432/postgres?schema=public"
   go run main.go server
   ```

## Test Methods

### Method 1: Health Check Endpoint

Check if RabbitMQ is connected via the health endpoint:

```bash
curl http://localhost:6906/readyz
```

Expected response (if RabbitMQ is working):
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
  },
  "timestamp": "2025-12-23T..."
}
```

### Method 2: RabbitMQ Test Endpoint

Test the RabbitMQ connection directly:

```bash
curl http://localhost:6906/rabbitmq/test
```

Expected response:
```json
{
  "status": "ok",
  "message": "RabbitMQ connection is working",
  "config": {
    "host": "localhost",
    "port": "5672",
    "enabled": true,
    "username": "admin"
  }
}
```

### Method 3: Setup Test Exchange and Queue

Set up a test exchange and queue:

```bash
curl -X POST http://localhost:6906/rabbitmq/setup
```

Expected response:
```json
{
  "status": "ok",
  "message": "Exchange and queue setup successfully",
  "exchange": "test-events",
  "queue": "test-queue",
  "routing_key": "test.*"
}
```

### Method 4: Publish a Test Message

Publish a test message to RabbitMQ:

```bash
curl -X POST http://localhost:6906/rabbitmq/publish \
  -H "Content-Type: application/json" \
  -d '{
    "exchange": "test-events",
    "routingKey": "test.message",
    "message": {
      "test": true,
      "data": "Hello RabbitMQ!"
    }
  }'
```

Expected response:
```json
{
  "status": "ok",
  "message": "Message published successfully",
  "event_id": "550e8400-e29b-41d4-a716-446655440000",
  "exchange": "test-events",
  "routing_key": "test.message",
  "published_at": "2025-12-23T..."
}
```

### Method 5: Verify Message in RabbitMQ Management UI

1. Open RabbitMQ Management UI: http://localhost:15672
2. Login with credentials (default: `admin` / `admin`)
3. Navigate to **Queues** tab
4. You should see `test-queue` with 1 message
5. Click on the queue to see message details
6. Click **Get messages** to retrieve the message

### Method 6: Test with a Consumer (Programmatic)

Create a simple Go test file to consume messages:

```go
package main

import (
	"fmt"
	"log"
	"time"
	
	"github.com/oskargbc/dws-event-service.git/internal/pkg/rabbitmq"
)

func main() {
	// Setup exchange and queue first
	err := rabbitmq.SetupExchangeAndQueue(
		"test-events",
		"topic",
		"test-queue",
		"test.*",
		true,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Start consumer
	handler := func(msg *rabbitmq.EventMessage) error {
		fmt.Printf("Received event: %s\n", msg.EventType)
		fmt.Printf("Event ID: %s\n", msg.EventID)
		fmt.Printf("Data: %+v\n", msg.Data)
		fmt.Printf("Timestamp: %s\n", msg.Timestamp)
		return nil
	}

	err = rabbitmq.StartEventConsumer("test-queue", "test-consumer", handler)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Consumer started. Waiting for messages...")
	
	// Keep running
	select {}
}
```

## Complete Test Workflow

Here's a complete workflow to test everything:

```bash
# Terminal 1: Port-forward RabbitMQ
kubectl port-forward -n rabbitmq svc/rabbitmq-service 5672:5672

# Terminal 2: Start your service
export DATABASE_URL="postgresql://postgres:w5f7eE-CJrTjPl-DfYZkV@localhost:5432/postgres?schema=public"
go run main.go server

# Terminal 3: Run tests
# 1. Test connection
curl http://localhost:6906/rabbitmq/test

# 2. Setup exchange and queue
curl -X POST http://localhost:6906/rabbitmq/setup

# 3. Publish a message
curl -X POST http://localhost:6906/rabbitmq/publish \
  -H "Content-Type: application/json" \
  -d '{
    "exchange": "test-events",
    "routingKey": "test.message",
    "message": {
      "test": true,
      "data": "Hello from curl!"
    }
  }'

# 4. Check health endpoint
curl http://localhost:6906/readyz | jq '.checks.rabbitmq'
```

## Troubleshooting

### Connection Failed

If you get connection errors:

1. **Check port-forward is running**:
   ```bash
   # Should show the port-forward process
   ps aux | grep port-forward
   ```

2. **Verify RabbitMQ service exists**:
   ```bash
   kubectl get svc -n rabbitmq
   ```

3. **Check credentials** in `configs/config.yaml`:
   ```yaml
   rabbitmq:
     username: "admin"
     password: "admin"
   ```

### Message Not Appearing in Queue

1. Make sure you've run the setup endpoint first
2. Check the routing key matches the binding pattern
3. Verify the exchange type matches (topic, direct, etc.)
4. Check RabbitMQ Management UI for errors

### Health Check Fails

If `/readyz` shows RabbitMQ as failed:

1. Check service logs for connection errors
2. Verify RabbitMQ is enabled in config: `rabbitmq.enabled: true`
3. Check network connectivity to RabbitMQ

## Expected Logs

When RabbitMQ connects successfully, you should see:

```
{"level":"info","msg":"Connecting to RabbitMQ at localhost:5672","time":"..."}
{"level":"info","msg":"Successfully connected to RabbitMQ","time":"..."}
{"level":"info","msg":"RabbitMQ connection verified successfully","time":"..."}
```

## Next Steps

Once testing is complete, you can:

1. **Integrate RabbitMQ into your event creation flow** - publish events when they're created
2. **Set up consumers** for other services to consume events
3. **Configure dead letter queues** for failed message handling
4. **Set up monitoring** and alerting for RabbitMQ

