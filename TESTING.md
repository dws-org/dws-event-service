# Testing Documentation

## Overview

This document describes the testing strategy and coverage for the DWS Event Service backend.

## Requirements Fulfilled

### REQ5: Code Coverage
- **Target:** Minimum 50% code coverage on backend services
- **Status:** ✅ Implemented
- **Coverage Areas:**
  - Controllers (events, health, rabbitmq)
  - Middlewares (authentication, error handling)
  - Services (database, RabbitMQ)

### REQ7: Endpoint Failure Test Cases
- **Target:** At least 2 failure test cases per service/component
- **Status:** ✅ Implemented

#### Events Controller Failure Tests
1. **404 Not Found**: `TestGetEventByID_NotFound` - Tests event retrieval with non-existent ID
2. **400 Bad Request (Missing Fields)**: `TestCreateEvent_BadRequest_MissingFields` - Tests event creation with incomplete payload
3. **400 Bad Request (Invalid JSON)**: `TestCreateEvent_BadRequest_InvalidJSON` - Tests event creation with malformed JSON
4. **400 Bad Request (Invalid Data Types)**: `TestCreateEvent_BadRequest_InvalidDataTypes` - Tests event creation with wrong data types

#### Authentication Middleware Failure Tests
1. **401 Unauthorized (Missing Token)**: `TestKeycloakAuthMiddleware_MissingToken` - Tests protected endpoint without token
2. **401 Unauthorized (Invalid Format)**: `TestKeycloakAuthMiddleware_InvalidTokenFormat` - Tests with malformed Authorization header
3. **401 Unauthorized (Empty Token)**: `TestKeycloakAuthMiddleware_EmptyToken` - Tests with empty Bearer token

## Running Tests

### Run All Tests
```bash
go test ./internal/... -v
```

### Run Tests with Coverage
```bash
go test ./internal/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Run Specific Package Tests
```bash
# Events controller tests
go test ./internal/controllers/events -v

# Middleware tests
go test ./internal/middlewares -v

# Service tests
go test ./internal/services -v
```

### Check Coverage Percentage
```bash
go test ./internal/... -coverprofile=coverage.out
go tool cover -func=coverage.out | grep total
```

## Test Structure

```
internal/
├── controllers/
│   ├── events/
│   │   ├── controller.go
│   │   └── controller_test.go      # ✅ 4+ failure tests
│   └── health/
│       ├── controller.go
│       └── controller_test.go      # ✅ 3 tests
├── middlewares/
│   ├── auth.go
│   ├── keycloak_auth.go
│   └── auth_test.go                # ✅ 3 failure tests (401)
└── services/
    ├── database.go
    └── database_test.go            # ✅ Unit tests
```

## Test Categories

### 1. Unit Tests
- Test individual functions and methods
- Mock external dependencies
- Fast execution
- Examples: `TestNewController`, `TestGetDatabaseServiceInstance`

### 2. Integration Tests
- Test HTTP endpoints end-to-end
- Use httptest for request/response testing
- Examples: `TestGetEvents_Success`, `TestCreateEvent_BadRequest`

### 3. Failure Tests (REQ7)
- **401 Unauthorized**: Authentication failures
  - Missing token
  - Invalid token format
  - Expired token
- **400 Bad Request**: Validation failures
  - Missing required fields
  - Invalid JSON
  - Wrong data types
- **404 Not Found**: Resource not found
  - Non-existent IDs
- **500 Internal Server Error**: Database/service failures

## CI/CD Integration

Tests are automatically run in GitHub Actions on:
- Every push to `main`
- Every pull request

### CI Pipeline Steps
1. **Test Job**: Runs all tests with coverage
   - Checks minimum 50% coverage requirement
   - Uploads coverage to Codecov
   - Generates HTML coverage report
2. **Lint Job**: Runs golangci-lint for code quality
3. **Build Job**: Only runs if tests pass

## Coverage Requirements

| Package | Target Coverage | Status |
|---------|----------------|--------|
| Controllers | 60%+ | ✅ |
| Middlewares | 70%+ | ✅ |
| Services | 50%+ | ✅ |
| **Overall** | **50%+** | ✅ |

## Best Practices

1. **Test Naming**: Use descriptive names with pattern `Test<Function>_<Scenario>`
2. **Setup/Teardown**: Use `init()` for test mode setup
3. **Assertions**: Use `testify/assert` for readable assertions
4. **Mocking**: Mock database and external services for unit tests
5. **Coverage**: Aim for >50% overall, focus on critical paths

## Security Testing

### SQL Injection Protection (REQ22)
- Prisma ORM provides parameterized queries by default
- No raw SQL queries used
- All user inputs are sanitized through Go struct binding

### XSS Protection (REQ22)
- JSON responses are automatically escaped by Gin
- No HTML rendering on backend
- Frontend handles sanitization

## Future Improvements

- [ ] Add performance tests with benchmarking
- [ ] Add end-to-end tests with real database
- [ ] Increase coverage to 80%+
- [ ] Add mutation testing
- [ ] Add contract testing for API
