# Go Microservice Template

An opinionated Go microservice boilerplate. Clone it once per service, run a single setup script, and you are ready to add domain logic.

## Highlights

- **Per-service metadata** baked into the CLI, logging, and config.
- **Setup automation** via `scripts/setup-service.sh` for renaming module paths and service identifiers.
- **Production-friendly HTTP** server with graceful shutdown and structured logging.
- **Health + metadata endpoints** (`/livez`, `/readyz`, `/healthz`, `/_meta`) for Kubernetes or any orchestrator.
- **Auth foundations** with JWT middleware, plus Supabase client wiring.
- **Prisma ORM** integration ready for Postgres (or any Prisma-supported database).

## Repository Layout

```
.
├── cmd/                 # Cobra CLI (root + server commands)
├── configs/             # YAML configuration + strongly typed structs
├── internal/            # Service implementation (router, controllers, etc.)
│   ├── controllers/
│   │   └── health/      # Reusable liveness/readiness handlers
│   ├── middlewares/
│   ├── pkg/             # Logger, Supabase client, utilities
│   ├── router/          # HTTP router setup & route registration
│   └── services/        # Infrastructure services (database, prisma wrapper)
├── prisma/              # Prisma schema + generated client
├── scripts/             # Tooling (service setup helper)
├── Dockerfile
├── Makefile
└── main.go
```

## One-Time Setup Per Service

1. Clone this template (or click "Use this template" on GitHub).
2. Run the setup helper to rename the module, service metadata, and port:
   ```bash
   ./scripts/setup-service.sh github.com/acme/user-service "User Service" 8081
   ```
3. Review and commit the generated changes, then push to the new repository.

The script updates:

- `go.mod` module path + all Go import statements.
- `configs/config.yaml` service name, slug, and default port.
- CLI help/output (e.g. `service version` command).

Repeat this workflow for each microservice you need (User, Chat, Event, Ticket, ...).

## Local Development

```bash
# Install dependencies
go mod download

# Launch the API
go run main.go server

# Run with hot reload (requires air)
air
```

The server listens on the port defined in `configs/config.yaml` (`service.server.port`), defaulting to `:8080`.

### Useful Commands

- `go test ./...` — run unit tests
- `make build` — produce a production binary
- `docker build -t my-service .` — build a container image

## Configuration

`configs/config.yaml` is the single source of truth. Notable sections:

- `service`: name, slug, description, version, tags.
- `server`: `host`, `port`, `gin_mode`.
- `jwt`: secret placeholder for auth middleware.
- `supabase`: base wiring if you integrate with Supabase.

Override values via environment variables or additional config files if desired.

## HTTP Surface

- `GET /livez` — liveness probe.
- `GET /readyz` / `GET /healthz` — readiness probe (extend `DatabaseService.HealthCheck` for real dependency checks).
- `GET /_meta` — service metadata (name, version, uptime, tags).

Add feature routes inside `internal/controllers` and register them in `internal/router/router.go`.

## Database & Prisma

1. Add models to `prisma/schema.prisma`.
2. Regenerate the client:
   ```bash
   cd prisma
   prisma generate
   prisma db push
   ```
3. Extend `internal/services/database.go` or create domain services that leverage the generated Prisma client.

`DatabaseService.HealthCheck` is intentionally light — plug in domain-specific checks (e.g. `SELECT 1`) once your schema exists.

## Supabase & Authentication

- `internal/middlewares/jwt.go` and `internal/pkg/utils/jwt.go` include helper utilities for JWT-based auth.
- `internal/pkg/supabase` instantiates the Supabase client with values from config/env.

Replace or expand these components per microservice requirements.

## Microservices Blueprint

For the architecture outlined in `docs/project-proposal-group6.md`, treat this repository as the base for each service (User, Chat, Event, Ticket, ...). Each clone yields an isolated Go module with consistent operational tooling, making it easy to deploy independently while preserving cross-cutting conventions.
