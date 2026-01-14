#!/bin/bash

# Test only packages that work without database
go test \
  ./configs \
  ./internal/controllers \
  ./internal/controllers/events \
  ./internal/controllers/health \
  ./internal/event \
  ./internal/middlewares \
  ./internal/pkg/logger \
  ./internal/pkg/metrics \
  ./internal/pkg/utils \
  ./internal/router \
  -coverprofile=coverage_working.out -covermode=atomic

# Calculate average coverage
echo "=== Package Coverage ==="
go tool cover -func=coverage_working.out | grep -E "\.go:" | awk '{print $3}' | sed 's/%//' | awk '{sum+=$1; count++} END {print "Average Coverage:", sum/count "%"}'

echo ""
echo "=== Total Coverage ==="
go tool cover -func=coverage_working.out | grep total
