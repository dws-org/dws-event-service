# Build stage
FROM golang:1.23.1 AS builder

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Install swag for generating swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy source code
COPY . .

# Generate Swagger docs
RUN swag init -g cmd/root.go --output ./docs

# Generate Prisma client
RUN go run github.com/steebchen/prisma-client-go generate

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Runtime stage - use Debian-based image for Prisma compatibility
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates openssl && rm -rf /var/lib/apt/lists/*

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs
# Copy Prisma binaries
COPY --from=builder /tmp/prisma /tmp/prisma

# Port freigeben
EXPOSE 6906

# Run the binary
CMD ["./main"]
