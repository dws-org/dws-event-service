#!/bin/bash

# RabbitMQ Test Script
# This script tests RabbitMQ connectivity and functionality

BASE_URL="http://localhost:6906"

echo "🧪 Testing RabbitMQ Integration"
echo "================================"
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test 1: Health Check
echo -e "${YELLOW}Test 1: Health Check${NC}"
echo "GET $BASE_URL/readyz"
response=$(curl -s "$BASE_URL/readyz")
rabbitmq_status=$(echo "$response" | grep -o '"rabbitmq":[^}]*' | grep -o '"status":"[^"]*"' | cut -d'"' -f4)

if [ "$rabbitmq_status" = "ok" ]; then
    echo -e "${GREEN}✓ RabbitMQ is healthy${NC}"
else
    echo -e "${RED}✗ RabbitMQ health check failed${NC}"
    echo "Response: $response"
fi
echo ""

# Test 2: Connection Test
echo -e "${YELLOW}Test 2: Connection Test${NC}"
echo "GET $BASE_URL/rabbitmq/test"
response=$(curl -s "$BASE_URL/rabbitmq/test")
status=$(echo "$response" | grep -o '"status":"[^"]*"' | cut -d'"' -f4)

if [ "$status" = "ok" ]; then
    echo -e "${GREEN}✓ RabbitMQ connection is working${NC}"
    echo "$response" | python3 -m json.tool 2>/dev/null || echo "$response"
else
    echo -e "${RED}✗ RabbitMQ connection failed${NC}"
    echo "$response"
fi
echo ""

# Test 3: Setup Exchange and Queue
echo -e "${YELLOW}Test 3: Setup Exchange and Queue${NC}"
echo "POST $BASE_URL/rabbitmq/setup"
response=$(curl -s -X POST "$BASE_URL/rabbitmq/setup")
status=$(echo "$response" | grep -o '"status":"[^"]*"' | cut -d'"' -f4)

if [ "$status" = "ok" ]; then
    echo -e "${GREEN}✓ Exchange and queue setup successful${NC}"
    echo "$response" | python3 -m json.tool 2>/dev/null || echo "$response"
else
    echo -e "${RED}✗ Setup failed${NC}"
    echo "$response"
fi
echo ""

# Test 4: Publish Test Message
echo -e "${YELLOW}Test 4: Publish Test Message${NC}"
echo "POST $BASE_URL/rabbitmq/publish"
response=$(curl -s -X POST "$BASE_URL/rabbitmq/publish" \
  -H "Content-Type: application/json" \
  -d '{
    "exchange": "test-events",
    "routingKey": "test.message",
    "message": {
      "test": true,
      "data": "Hello RabbitMQ!",
      "timestamp": "'$(date -u +"%Y-%m-%dT%H:%M:%SZ")'"
    }
  }')
status=$(echo "$response" | grep -o '"status":"[^"]*"' | cut -d'"' -f4)

if [ "$status" = "ok" ]; then
    echo -e "${GREEN}✓ Message published successfully${NC}"
    echo "$response" | python3 -m json.tool 2>/dev/null || echo "$response"
else
    echo -e "${RED}✗ Failed to publish message${NC}"
    echo "$response"
fi
echo ""

echo "================================"
echo -e "${GREEN}Testing complete!${NC}"
echo ""
echo "💡 Tip: Check RabbitMQ Management UI at http://localhost:15672"
echo "   to see the published message in the test-queue"

