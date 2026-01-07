# Event Service API Documentation

## Overview

The Event Service is responsible for managing events in the DWS platform. It provides endpoints for:
- Listing all available events
- Viewing detailed event information
- Creating new events (organizers only)

## Base URLs

- **Production**: `https://event.ltu-m7011e-6.se`
- **Swagger UI**: `https://event.ltu-m7011e-6.se/swagger/index.html`

## Authentication

All `/api/v1/*` endpoints require Bearer token authentication.

### Getting a Token

1. Navigate to: `https://keycloak.ltu-m7011e-6.se/realms/dws-org`
2. Login with your credentials
3. Use the JWT token in API requests:

```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  https://event.ltu-m7011e-6.se/api/v1/events
```

## Endpoints

### GET /api/v1/events

List all events.

**Authentication**: Required  
**Authorization**: All authenticated users

**Response**: `200 OK`
```json
[
  {
    "id": "evt-001",
    "name": "Rock Festival 2026",
    "description": "Annual rock music festival",
    "startDate": "2026-06-15T00:00:00Z",
    "startTime": "2026-06-15T18:00:00Z",
    "endDate": "2026-06-17T00:00:00Z",
    "location": "Stockholm Arena",
    "capacity": 5000,
    "price": 599.00,
    "imageUrl": "https://example.com/festival.jpg",
    "category": "Music",
    "organizerId": "org-123",
    "createdAt": "2026-01-01T10:00:00Z",
    "updatedAt": "2026-01-01T10:00:00Z"
  }
]
```

### GET /api/v1/events/{id}

Get a single event by ID.

**Authentication**: Required  
**Authorization**: All authenticated users

**Parameters**:
- `id` (path) - Event ID

**Response**: `200 OK`
```json
{
  "id": "evt-001",
  "name": "Rock Festival 2026",
  ...
}
```

**Error Responses**:
- `401 Unauthorized` - Missing or invalid token
- `404 Not Found` - Event does not exist

### POST /api/v1/events

Create a new event.

**Authentication**: Required  
**Authorization**: Users with `Organiser` realm role

**Request Body**:
```json
{
  "name": "Jazz Night",
  "description": "Evening of smooth jazz",
  "startDate": "2026-08-20T00:00:00Z",
  "startTime": "2026-08-20T19:00:00Z",
  "endDate": "2026-08-20T00:00:00Z",
  "location": "Blue Note Club",
  "capacity": 200,
  "price": 299.00,
  "imageUrl": "https://example.com/jazz.jpg",
  "category": "Music",
  "organizerId": "org-456"
}
```

**Response**: `201 Created`
```json
{
  "id": "evt-002",
  "name": "Jazz Night",
  ...
}
```

**Error Responses**:
- `400 Bad Request` - Invalid request payload
- `401 Unauthorized` - Missing or invalid token
- `403 Forbidden` - User doesn't have Organiser role
- `500 Internal Server Error` - Failed to create event

## Health Checks

### GET /livez

Kubernetes liveness probe.

**Response**: `200 OK`
```json
{
  "status": "alive"
}
```

### GET /readyz

Kubernetes readiness probe.

**Response**: `200 OK`
```json
{
  "status": "ready"
}
```

## Error Handling

All error responses follow this format:

```json
{
  "error": "error_code",
  "message": "Human-readable message",
  "details": "Additional details"
}
```

Common error codes:
- `unauthorized` - Authentication required
- `forbidden` - Insufficient permissions
- `not_found` - Resource not found
- `invalid_request` - Bad request payload
- `internal_error` - Server error

## Rate Limiting

Currently no rate limiting is implemented.

## Versioning

API version is included in the URL path: `/api/v1/`

## Support

For issues or questions:
- GitHub: https://github.com/dws-org/dws-event-service
- OpenAPI Spec: `docs/openapi.yaml`
